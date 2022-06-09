package slice

func Map[A any, B any](in []A, f func(A) B) []B {
	res := make([]B, 0, len(in))
	for _, el := range in {
		res = append(res, f(el))
	}
	return res
}
