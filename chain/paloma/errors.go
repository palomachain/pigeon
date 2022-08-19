package paloma

import (
	"errors"
	"net"

	"github.com/vizualni/whoops"
)

const (
	ErrUnableToDecodeAddress = whoops.Errorf("unable to decode address: %s")
	ErrNodeIsNotInSync       = whoops.String("paloma node is not in sync")

	ErrPalomaIsDown = whoops.String("paloma is down")
)

func IsPalomaDown(err error) bool {
	var netErr *net.OpError
	return errors.As(err, &netErr)
}
