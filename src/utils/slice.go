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

func SliceIncluded[T comparable](slice []T, val T) bool {
	bl := false
	for _, item := range slice {
		if item == val {
			bl = true
			break
		} else {
			bl = false
		}
	}
	return bl
}
func SliceIncludeds[T comparable](slice []T, vals []T) bool {
	bl := false
	for _, val := range vals {
		bl = SliceIncluded(slice, val)
		if bl {
			break
		}
	}
	return bl
}
