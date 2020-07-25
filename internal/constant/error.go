package constant

import "errors"

var (
	// ErrRequiredDB will throw if database address is empty.
	ErrRequiredDB = errors.New("required database address")
	// ErrInvalidDB will throw if database address format is invalid.
	ErrInvalidDB = errors.New("invalid database address")
	// ErrRequiredRabbit will throw if rabbitmq connection is empty.
	ErrRequiredRabbit = errors.New("required rabbitmq connection")
	// ErrNotFound will throw if data not found.
	ErrNotFound = errors.New("data not found")
	// ErrRequiredUser will throw if user id is empty.
	ErrRequiredUser = errors.New("required user id")
	// ErrInvalidPoint will throw if point is zero or lower.
	ErrInvalidPoint = errors.New("point must be positive and not zero")
)
