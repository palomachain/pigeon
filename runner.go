package runner

import (
	"context"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/strangelove-ventures/lens/byop"
	lens "github.com/strangelove-ventures/lens/client"
	"github.com/volumefi/cronchain/runner/client/cronchain"
	cronchaintypes "github.com/volumefi/cronchain/runner/client/cronchain/types"
	"github.com/volumefi/cronchain/runner/client/terra/types"

	chain "github.com/volumefi/cronchain/runner/client"
)

func Start() {
	panic("main loop goes here")
	// var c terra.Client
	// modules := append(lens.ModuleBasics[:], byop.NewModule(
	// 	"testing",
	// 	(*types.MsgExecuteContract)(nil),
	// ))

	// modules := append(lens.ModuleBasics[:])
	// lensc, err := chain.NewChainClient(
	// 	&lens.ChainClientConfig{
	// 		Key:            "matija",
	// 		ChainID:        "columbus-5",
	// 		RPCAddr:        "http://127.0.0.1:22222",
	// 		KeyringBackend: "file",
	// 		KeyDirectory:   "/home/vizualni/.terra/",
	// 		AccountPrefix:  "terra",
	// 		Modules:        modules,
	// 		Debug:          true,
	// 		GasAdjustment:  1.1,
	// 		GasPrices:      "0.2056735uusd,0luna",
	// 		SignModeStr:    "direct",
	// 	},
	// 	"doesn-tmatter",
	// 	os.Stdin,
	// 	os.Stdout,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// c.LensClient = lensc
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// c.ExecuteSmartContract(ctx)
}

func Start2() {
	var c cronchain.Client
	byopmodule := byop.NewModule(
		"testing",
		(*types.MsgExecuteContract)(nil),
	)
	byopmodule.MsgsInterfaces = []byop.RegisterInterface{
		{
			Name:  "volumefi.cronchain.concensus.QueuedSignedMessage",
			Iface: (*cronchaintypes.QueuedSignedMessageI)(nil),
			Msgs: []proto.Message{
				&cronchaintypes.QueuedSignedMessage{},
			},
		},
	}
	byopmodule.MsgsImplementations = []byop.RegisterImplementation{
		{
			Iface: (*cronchaintypes.QueuedSignedMessageI)(nil),
			Msgs: []proto.Message{
				&cronchaintypes.QueuedSignedMessage{},
			},
		},
	}
	byopmodule.MsgsAmino = []byop.RegisterAmino{
		{
			ConcreteIface: &cronchaintypes.QueuedSignedMessage{},
			Name:          "concensus/QueuedSignedMessage",
		},
		{
			ConcreteIface: &cronchaintypes.Signer{},
			Name:          "concensus/Signer",
		},
	}
	modules := append(lens.ModuleBasics[:], byopmodule)

	lensc, err := chain.NewChainClient(
		&lens.ChainClientConfig{
			Key:            "matija",
			ChainID:        "cronchain",
			RPCAddr:        "http://127.0.0.1:26657",
			KeyringBackend: "test",
			KeyDirectory:   "/home/vizualni/.cronchain/",
			AccountPrefix:  "stake",
			Modules:        modules,
			Debug:          true,
			GasAdjustment:  1.1,
			GasPrices:      "0.2056735uusd,0luna",
			SignModeStr:    "direct",
		},
		"doesn-tmatter",
		os.Stdin,
		os.Stdout,
	)
	if err != nil {
		panic(err)
	}
	c.L = lensc
	c.QueryMessagesForExecution(context.Background())

}
