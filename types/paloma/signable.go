package paloma

type Signable interface {
	Signable()
}

func (SignSmartContractExecute) Signable() {}
