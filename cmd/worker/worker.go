package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rl404/point-system/internal/config"
	"github.com/rl404/point-system/internal/controller"
	"github.com/rl404/point-system/internal/model"
)

// DB connection.
var DB *gorm.DB

func main() {
	// Get config.
	cfg := config.GetConfig()

	// Init db connection.
	var err error
	DB, err = cfg.InitDB()
	if err != nil {
		log.Fatal("error init db", " - ", err.Error())
		return
	}

	// Init rabbitmq connection.
	rabbit, err := cfg.InitRabbit()
	if err != nil {
		log.Fatal("error init rabbitmq", " - ", err.Error())
		return
	}

	// Consume queue.
	msgs, err := rabbit.Consume(
		config.RabbitMQName, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)

	if err != nil {
		log.Fatal("error consume rabbitmq", "-", err.Error())
		return
	}

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			var data controller.RabbitQueue
			json.Unmarshal(m.Body, &data)

			err := updatePoint(data)
			if err != nil {
				fmt.Printf("failed updating point - %s \n", err.Error())
			}

			m.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// updatePoint to update user point.
func updatePoint(data controller.RabbitQueue) error {
	var userPoint model.UserPoint
	DB.Where("user_id = ?", data.UserID).First(&userPoint)

	userPoint.UserID = data.UserID
	if data.Action == controller.ActionAdd {
		userPoint.Point++
	} else if data.Action == controller.ActionSubstract {
		userPoint.Point--
	}

	err := DB.Save(&userPoint).Error
	if err != nil {
		return err
	}

	// Print notif.
	fmt.Printf("[%s] user %v - %s \n", time.Now().Format("2006-01-02 15:04:05.00"), userPoint.UserID, data.Action)

	return insertLog(data)
}

// insertLog to insert log.
func insertLog(data controller.RabbitQueue) error {
	log := model.Log{
		UserID:      data.UserID,
		Action:      data.Action,
		RequestedAt: data.RequestedAt,
	}
	return DB.Create(&log).Error
}
