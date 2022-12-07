package util

func Filter[S ~[]T, T any](s S, f func(T) bool) S {
	ret := make(S, 0, len(s))
	for _, e := range s {
		if f(e) {
			ret = append(ret, e)
		}
	}
	return ret
}

func Select[T, V any](s []T, f func(T) V) []V {
	ret := make([]V, 0, len(s))
	for _, e := range s {
		ret = append(ret, f(e))
	}
	return ret
}
