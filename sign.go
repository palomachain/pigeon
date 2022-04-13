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

func jsonDeterministicEncoding(msg any) ([]byte, error) {
	// take anything and create a json byte slice
	js, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	var c any
	// then unmarshal this back to Go!
	err = json.Unmarshal(js, &c)
	if err != nil {
		return nil, err
	}
	// and finally, when calling the marshal on the new unmarshaled data
	// it will be sorted!
	js, err = json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
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
