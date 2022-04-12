package conductor

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type signer interface {
	Sign([]byte) ([]byte, error)
}

type serializer interface {
	DeterministicSerialize(any) ([]byte, error)
}

type signerFnc func([]byte) ([]byte, error)

func (fnc signerFnc) Sign(b []byte) ([]byte, error) {
	return fnc(b)
}

func keyringSigner(k keyring.Signer, uid string) signerFnc {
	return signerFnc(func(b []byte) ([]byte, error) {
		signedBytes, _, err := k.Sign(uid, b)
		return signedBytes, err
	})
}

type serializeFnc func(any) ([]byte, error)

func (fnc serializeFnc) DeterministicSerialize(msg any) ([]byte, error) {
	return fnc(msg)
}

// TODO: verify that this sorts keys and that it's deterministic!!!
// keys should be lowercase and sorted.
func jsonDeterministicEncoding(msg any) ([]byte, error) {
	return json.Marshal(msg)
}

func signBytes(s signer, ser serializer, msg any, nonce []byte) ([]byte, []byte, error) {
	encodedMsg, err := ser.DeterministicSerialize(msg)

	if err != nil {
		return nil, nil, err
	}

	// appending nonce to the end of the message that needs to be signed
	msgWithNonce := append(encodedMsg, nonce...)

	signedBytes, err := s.Sign(msgWithNonce)
	if err != nil {
		return nil, nil, err
	}

	return signedBytes, msgWithNonce, nil
}
