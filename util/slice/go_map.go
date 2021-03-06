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

// MakeMapKeys makes a map of provided slice and a function which
// returns a key value for a map given an item from a slice.
// If key already exists, it overrides it.
func MakeMapKeys[K comparable, V any](slice []V, getKey func(V) K) map[K]V {
	m := make(map[K]V, len(slice))
	for _, item := range slice {
		key := getKey(item)
		m[key] = item
	}
	return m
}
