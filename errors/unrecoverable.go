package errors

import (
	"errors"

	"github.com/vizualni/whoops"
)

const (
	ErrUnrecoverable = whoops.String("unrecoverable error")
)

func IsUnrecoverable(err error) bool {
	return errors.Is(err, ErrUnrecoverable)
}

func Unrecoverable(err error) error {
	return whoops.Wrap(err, ErrUnrecoverable)
}
