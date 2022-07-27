package slice

import (
	"fmt"
)

func IterN[T any](num int, fn func(index int) T) (res []T) {
	if num <= 0 {
		panic("num must be a positive integer")
	}
	for i := 0; i < num; i++ {
		res = append(res, fn(i))
	}
	return
}

func IterMapN[K comparable, V any](num int, fn func(index int) (K, V)) (res map[K]V) {
	if num <= 0 {
		panic("num must be a positive integer")
	}
	for i := 0; i < num; i++ {
		k, v := fn(i)
		if _, ok := res[k]; ok {
			panic(fmt.Sprintf("key %s was already calculated", (any)(k)))
		}
		res[k] = v
	}
	return
}
