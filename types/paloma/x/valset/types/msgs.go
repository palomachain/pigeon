package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/market module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)
