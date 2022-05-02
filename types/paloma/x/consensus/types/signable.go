package types

type Signable interface {
	Signable()
}

func (SignSmartContractExecute) Signable() {}
