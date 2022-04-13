package cronchain

import (
	"github.com/vizualni/whoops"
)

const (
	ErrUnableToDecodeAddress       = whoops.Errorf("unable to decode address: %s")
	ErrUnableToUnpackAny           = whoops.String("unable to unpack Any message")
	ErrIncorrectTypeSavedInMessage = whoops.Errorf("expected type: %T, got type %T")
)
