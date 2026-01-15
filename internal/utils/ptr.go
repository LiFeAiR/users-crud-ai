package utils

func Ptr[A any](a A) *A {
	return &a
}

func FromPtr[A any](a *A) A {
	var result A
	if a != nil {
		result = *a
	}
	return result
}
