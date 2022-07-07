package chain

import (
	"github.com/vizualni/whoops"
)

const (
	ErrSignatureVerificationFailed              = whoops.String("signature verification failed")
	ErrSignatureDoesNotMatchItsRespectiveSigner = whoops.String("signature does not match its respective signer")
	ErrTooLittleOrTooManySignaturesProvided     = whoops.String("too many or too little signatures provided")
	ErrProcessorDoesNotSupportThisQueue         = whoops.Errorf("processor does not support queue: %s")

	ErrNotFound = whoops.String("not found")

	EnrichedChainReferenceID whoops.Field[string] = "chainReferenceID"
	EnrichedID               whoops.Field[uint64] = "id"
	EnrichedItemType         whoops.Field[string] = "type"
)
