package util

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func Reduce[T, A any](ts []T, acc A, fn func(value T, acc A) A) A {
	if len(ts) == 0 {
		return acc
	}

	return Reduce(ts[1:], fn(ts[0], acc), fn)
}
