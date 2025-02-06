package utils

// 查找项
func Find[T comparable](slice []T, fn func(v T, i int) bool) T {
	for i, val := range slice {
		if fn(val, i) {
			return val
		}
	}
	return *new(T)
}
