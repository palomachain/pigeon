package ion

type _err string

func (e _err) Error() string { return string(e) }

const (
	ErrTimeoutAfterWaitingForTxBroadcast _err = "timed out after waiting for tx to get included in the block"
	ErrUnexpectedNonZeroCode             _err = "node returned unexpected code"
)
