package byop

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/gogoproto/proto"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
)

var _ module.AppModuleBasic = Module{}

type RegisterInterface struct {
	Name  string
	Iface interface{}
	Msgs  []proto.Message
}

type RegisterImplementation struct {
	Iface interface{}
	Msgs  []proto.Message
}

type Module struct {
	ModuleName string

	MsgsInterfaces      []RegisterInterface
	MsgsImplementations []RegisterImplementation
}

// RegisterInterfaces is the only method that we care about. It registers the
// injected interfaces into the provided registry, so that it can be decoded.
func (m Module) RegisterInterfaces(registry types.InterfaceRegistry) {
	for _, mi := range m.MsgsInterfaces {
		registry.RegisterInterface(mi.Name, mi.Iface, mi.Msgs...)
	}
	for _, mi := range m.MsgsImplementations {
		registry.RegisterImplementations(mi.Iface, mi.Msgs...)
	}
}

// All other methods below exist just to fulfill the module.AppModuleBasic interface.

func (m Module) Name() string { return m.ModuleName }

func (m Module) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
}

func (m Module) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	panic("not required")
}

func (m Module) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error {
	panic("not required")
}

func (m Module) RegisterRESTRoutes(client.Context, *mux.Router) { panic("not required") }

func (m Module) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) { panic("not required") }

func (m Module) GetTxCmd() *cobra.Command { panic("not required") }

func (m Module) GetQueryCmd() *cobra.Command { panic("not required") }
