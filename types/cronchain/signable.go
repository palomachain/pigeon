package cronchain

type Signable interface {
	Signable()
}

func (SignSmartContractExecute) Signable() {}
