package slice

func Filter[V any](slice []V, filter func(el V) bool) []V {
	ret := []V{}
	for _, item := range slice {
		if filter(item) {
			ret = append(ret, item)
		}
	}
	return ret
}
