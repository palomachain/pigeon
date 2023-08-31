package libchain

import "math/big"

const cArbitrumID int64 = 42161

func IsArbitrum(id *big.Int) bool {
	return id.Int64() == cArbitrumID
}
