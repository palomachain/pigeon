package ion

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	log "github.com/sirupsen/logrus"
)

// This file is cursed and this mutex is too
// you don't want none of this dewey cox
var sdkConfigMutex sync.Mutex

// SetSDKConfig sets the SDK config to the proper bech32 prefixes
// Don't use this unless you know what you're doing.
// if lens is successful, this can be eliminated
// TODO: :dagger: :knife: :chainsaw: remove this function
func (cc *Client) SetSDKContext() func() {
	log.Info("Setting SDK context")
	sdkConfigMutex.Lock()
	log.Info("Setting config values")
	sdkConf := sdk.GetConfig()
	sdkConf.SetBech32PrefixForAccount(cc.Config.AccountPrefix, cc.Config.AccountPrefix+"pub")
	sdkConf.SetBech32PrefixForValidator(cc.Config.AccountPrefix+"valoper", cc.Config.AccountPrefix+"valoperpub")
	sdkConf.SetBech32PrefixForConsensusNode(cc.Config.AccountPrefix+"valcons", cc.Config.AccountPrefix+"valconspub")
	log.Info("returning the unlock")
	return sdkConfigMutex.Unlock
}
