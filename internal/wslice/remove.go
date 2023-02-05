package wslice

func FilterOne[T any](s []T, f func(T) bool) (*T, int) {
	for i, t := range s {
		if f(t) {
			return &t, i
		}
	}
	return nil, -1
}

func Remove[T any](s []T, r int) []T {
	if r >= 0 {
		return append(s[:r], s[r+1:]...)
	}

	return s
}
