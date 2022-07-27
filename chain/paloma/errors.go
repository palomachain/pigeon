package paloma

import (
	"github.com/vizualni/whoops"
)

const (
	ErrUnableToDecodeAddress = whoops.Errorf("unable to decode address: %s")
	ErrNodeIsNotInSync       = whoops.String("paloma node is not in sync")
)
