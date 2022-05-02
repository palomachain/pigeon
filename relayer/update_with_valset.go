package relayer

import (
	"context"

	"github.com/palomachain/sparrow/client/paloma"
	"github.com/palomachain/sparrow/errors"
	"github.com/vizualni/whoops"
)

func (r *Relayer) updateValidatorInfo(ctx context.Context) error {
	return whoops.Try(func() {
		whoops.Assert(
			r.registerValidator(ctx),
		)
		whoops.Assert(
			r.updateExternalChainInfos(ctx, "terra", r.config.Terra.Accounts),
		)
	})
}

func (r *Relayer) updateExternalChainInfos(ctx context.Context, chainID string, accAddresses []string) error {
	val, err := r.palomaClient.QueryValidatorInfo(ctx, r.validatorAddress)
	if err != nil {
		return err
	}
	chainInfos := []paloma.ChainInfoIn{}

	// TODO: implement accounts removal
	for _, accAddr := range accAddresses {
		// check if this acc is already registered
		found := false
		for _, currentChainInfo := range val.ExternalChainInfos {
			if currentChainInfo.ChainID == chainID && currentChainInfo.Address == accAddr {
				found = true
				break
			}
		}
		if !found {
			chainInfos = append(chainInfos, paloma.ChainInfoIn{
				ChainID:    chainID,
				AccAddress: accAddr,
			})

		}
	}

	if len(chainInfos) == 0 {
		return nil
	}

	return r.palomaClient.AddExternalChainInfo(ctx, chainInfos...)
}

func (r *Relayer) registerValidator(ctx context.Context) error {
	val, err := r.palomaClient.QueryValidatorInfo(ctx, r.validatorAddress)
	if err != nil {
		return nil
	}

	if val != nil {
		// already registered
		return nil
	}
	kr := r.palomaClient.Keyring()
	signingKeyName := r.config.Paloma.SigningKeyName
	keyInfo, err := kr.Key(signingKeyName)
	if err != nil {
		return errors.Unrecoverable(err)
	}

	pkBytes := keyInfo.GetPubKey().Bytes()

	sig, _, err := kr.Sign(signingKeyName, pkBytes)
	if err != nil {
		return errors.Unrecoverable(err)
	}

	return r.palomaClient.RegisterValidator(ctx, r.validatorAddress, pkBytes, sig)
}
