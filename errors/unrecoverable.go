package errors

import (
	"errors"

	"github.com/vizualni/whoops"
)

const (
	errUnrecoverable = whoops.String("unrecoverable error")
)

func IsUnrecoverable(err error) bool {
	return errors.Is(err, errUnrecoverable)
}

func Unrecoverable(err error) error {
	return whoops.Wrap(err, errUnrecoverable)
}
