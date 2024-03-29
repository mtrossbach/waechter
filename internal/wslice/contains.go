package wslice

func ContainsAll[T comparable](s []T, e []T) bool {
	for _, ee := range e {
		if !Contains(s, ee) {
			return false
		}
	}
	return true
}

func ContainsAny[T comparable](s []T, e []T) bool {
	for _, ee := range e {
		if Contains(s, ee) {
			return true
		}
	}
	return false
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
