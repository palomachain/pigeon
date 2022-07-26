package slice

func Filter[V any](slice []V, filters ...func(el V) bool) []V {
	ret := []V{}
	if len(filters) == 0 {
		return slice
	}
	for _, item := range slice {
		shouldAdd := true
		for _, filter := range filters {
			shouldAdd = shouldAdd && filter(item)
			if !shouldAdd {
				break
			}
		}
		if shouldAdd {
			ret = append(ret, item)
		}
	}
	return ret
}
