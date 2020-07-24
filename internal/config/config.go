package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rl404/point-system/internal/constant"
	"github.com/streadway/amqp"
)

// Config is config model for app.
type Config struct {
	// HTTP port.
	Port string

	// Database ip and port.
	Address string
	// Database name.
	DB string
	// Database schema name.
	Schema string
	// Database username.
	User string
	// Database password.
	Password string

	// RabbitMQ URL connection.
	RabbitmqUrl string `split_words:"true"`
}

const (
	// DefaultPort is default HTTP app port.
	DefaultPort = "31001"
	// EnvPath is .env file path.
	EnvPath = "../../config/config.env"
	// EnvPrefix is environment
	EnvPrefix = "PS"
	// DefaultMaxIdleConn is default max db idle connection.
	DefaultMaxIdleConn = 10
	// DefaultMaxOpenConn is default max db open connection.
	DefaultMaxOpenConn = 10
	// DefaultConnMaxLifeTime is default db connection max life time.
	DefaultConnMaxLifeTime = 5 * time.Minute
	// RabbitMQName is rabbitmq queue name.
	RabbitMQName = "POINT_SYSTEM_QUEUE"
)

// GetConfig to get config from env.
func GetConfig() (cfg Config) {
	cfg.Port = DefaultPort

	// Load .env file if exist.
	godotenv.Load(EnvPath)

	// Convert env to struct.
	envconfig.Process(EnvPrefix, &cfg)

	// Prepare the ":" for starting HTTP.
	cfg.Port = ":" + cfg.Port

	return cfg
}

// InitDB to intiate db connection.
func (c *Config) InitDB() (db *gorm.DB, err error) {
	if c.Address == "" {
		return nil, constant.ErrRequiredDB
	}

	// Split address and port.
	split := strings.Split(c.Address, ":")
	if len(split) != 2 {
		return nil, constant.ErrInvalidDB
	}

	// Open db connection.
	conn := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", split[0], split[1], c.DB, c.User, c.Password)
	db, err = gorm.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	// Set base connection setting.
	db.DB().SetMaxIdleConns(DefaultMaxIdleConn)
	db.DB().SetMaxOpenConns(DefaultMaxOpenConn)
	db.DB().SetConnMaxLifetime(DefaultConnMaxLifeTime)
	db.SingularTable(true)
	db.LogMode(false)

	// Set default schema.
	err = db.Exec(fmt.Sprintf("SET search_path TO %s", c.Schema)).Error
	if err != nil {
		return db, err
	}

	gorm.DefaultTableNameHandler = func(dbVeiculosGorm *gorm.DB, defaultTableName string) string {
		if c.Schema == "" {
			c.Schema = "public"
		}
		return c.Schema + "." + defaultTableName
	}

	// Validate db structure.
	err = c.validateDB(db)
	if err != nil {
		return db, err
	}

	return db, nil
}

// InitRabbit to intiate rabbitmq connection.
func (c *Config) InitRabbit() (ch *amqp.Channel, err error) {
	if c.RabbitmqUrl == "" {
		return ch, constant.ErrRequiredRabbit
	}

	conn, err := amqp.Dial(c.RabbitmqUrl)
	if err != nil {
		return ch, err
	}

	ch, err = conn.Channel()
	if err != nil {
		return ch, err
	}

	_, err = ch.QueueDeclare(
		RabbitMQName, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return ch, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	return ch, err
}
