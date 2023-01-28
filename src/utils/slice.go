package utils

type sliceType struct{}

var Slice sliceType

// 切片中是否包含
func (s *sliceType) Includes(slice []string, value string) bool {
	isRegister := false
	for _, val := range slice {
		if val == value {
			isRegister = true
			break
		}
	}
	return isRegister
}
