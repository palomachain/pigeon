package queue

import "strings"

const (
	QueueSuffixTurnstone          = "evm-turnstone-message"
	QueueSuffixValidatorsBalances = "validators-balances"
)

type TypeName string

func FromString(s string) TypeName {
	return TypeName(s)
}

func (t TypeName) String() string {
	return string(t)
}

func (t TypeName) IsTurnstoneQueue() bool {
	return strings.HasSuffix(string(t), QueueSuffixTurnstone)
}

func (t TypeName) IsValidatorsValancesQueue() bool {
	return strings.HasSuffix(string(t), QueueSuffixValidatorsBalances)
}
