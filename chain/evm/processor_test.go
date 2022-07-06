package evm

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/config"
	"github.com/stretchr/testify/require"
)

func TestEvmSigning(t *testing.T) {
	ks := OpenKeystore(t.TempDir())
	acc, err := ks.NewAccount("abcd")
	require.NoError(t, err)
	ks.Unlock(acc, "abcd")

	c := &Client{
		keystore: ks,
		addr:     acc.Address,
		config: config.EVM{
			ChainClientConfig: config.ChainClientConfig{
				SigningKey: acc.Address.Hex(),
			},
		},
	}

	p := Processor{
		evmClient: c,
	}
	require.NoError(t, err)

	msgsToSign := []chain.QueuedMessage{
		{
			BytesToSign: crypto.Keccak256([]byte("hello")),
		},
		{
			BytesToSign: crypto.Keccak256([]byte("world")),
		},
	}

	signed, err := p.SignMessages(context.Background(), "does not matter", msgsToSign...)
	require.NoError(t, err)
	for i := range msgsToSign {
		orig, signed := msgsToSign[i], signed[i]

		signedMsg := crypto.Keccak256(append([]byte(SignedMessagePrefix), orig.BytesToSign...))

		pubkey, err := crypto.Ecrecover(signedMsg, signed.Signature)
		require.NoError(t, err)

		pk, err := crypto.UnmarshalPubkey(pubkey)
		require.NoError(t, err)

		addr := crypto.PubkeyToAddress(*pk)

		require.Equal(t, acc.Address, addr)
	}
}

func TestProcessingMessages(t *testing.T) {
	ctx := context.Background()
	for _, tt := range []struct {
		name      string
		setup     func(t *testing.T) Processor
		queueName string
		msgs      []chain.MessageWithSignatures
		expErr    error
	}{
		{
			name:      "with unexpected queuename it returns an error",
			queueName: "i dont exist",
			setup: func(t *testing.T) Processor {
				return Processor{}
			},
			expErr: chain.ErrProcessorDoesNotSupportThisQueue,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			processor := tt.setup(t)
			err := processor.ProcessMessages(ctx, tt.queueName, tt.msgs)

			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
