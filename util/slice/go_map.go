package slice

func FromMapValues[K comparable, V any](mm map[K]V) []V {
	res := make([]V, 0, len(mm))
	for _, v := range mm {
		res = append(res, v)
	}
	return res
}

func FromMapKeys[K comparable, V any](mm map[K]V) []K {
	res := make([]K, 0, len(mm))
	for k := range mm {
		res = append(res, k)
	}
	return res
}
