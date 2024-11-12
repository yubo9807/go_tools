package m3

import (
	"github.com/tjfoc/gmsm/sm3"
)

// SM3 加密哈希计算
func Encrypt(data []byte) ([]byte, error) {
	// 使用 SM3 算法计算哈希值
	hash := sm3.New()
	hash.Write(data)

	// 返回哈希值
	return hash.Sum(nil), nil
}
