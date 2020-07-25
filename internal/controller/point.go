package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rl404/point-system/internal/config"
	"github.com/rl404/point-system/internal/constant"
	"github.com/rl404/point-system/internal/model"
	"github.com/rl404/point-system/internal/view"
	"github.com/streadway/amqp"
)

// PointHandler to handle all point activities.
type PointHandler struct {
	DB     *gorm.DB
	Rabbit *amqp.Channel
}

// newPointHandler to create new instance.
func newPointHandler(cfg config.Config) (ph PointHandler, err error) {
	// Init db connection.
	ph.DB, err = cfg.InitDB()
	if err != nil {
		return ph, err
	}

	// Init rabbitmq connection.
	ph.Rabbit, err = cfg.InitRabbit()
	if err != nil {
		return ph, err
	}

	return ph, nil
}

// getPoint to get user's point.
func (ph *PointHandler) getPoint(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		view.RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var userPoint model.UserPoint
	if ph.DB.Where("user_id = ?", id).First(&userPoint).RecordNotFound() {
		view.RespondWithJSON(w, http.StatusNotFound, constant.ErrNotFound.Error(), nil)
		return
	}

	view.RespondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK), userPoint)
}

// addPoint to send queue adding point to rabbitmq.
func (ph *PointHandler) addPoint(w http.ResponseWriter, r *http.Request) {
	var request RabbitQueue

	// Get request body.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		view.RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Validate request data.
	err = ph.validateRequest(request)
	if err != nil {
		view.RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Set action to add.
	request.Action = ActionAdd
	request.RequestedAt = time.Now()

	// Send to queue.
	err = ph.sendQueue(request)
	if err != nil {
		view.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	view.RespondWithJSON(w, http.StatusAccepted, http.StatusText(http.StatusAccepted), nil)
}

// subtractPoint to send queue subtracting point to rabbitmq.
func (ph *PointHandler) subtractPoint(w http.ResponseWriter, r *http.Request) {
	var request RabbitQueue

	// Get request body.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		view.RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Validate request data.
	err = ph.validateRequest(request)
	if err != nil {
		view.RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Set action to subtract.
	request.Action = ActionSubtract
	request.RequestedAt = time.Now()

	// Send to queue.
	err = ph.sendQueue(request)
	if err != nil {
		view.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	view.RespondWithJSON(w, http.StatusAccepted, http.StatusText(http.StatusAccepted), nil)
}

// validateRequest to validate request data.
func (ph *PointHandler) validateRequest(request RabbitQueue) error {
	if request.UserID == 0 {
		return constant.ErrRequiredUser
	}

	if request.Point <= 0 {
		return constant.ErrInvalidPoint
	}

	return nil
}
