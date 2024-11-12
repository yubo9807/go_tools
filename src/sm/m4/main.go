package m4

import (
	"bytes"
	"crypto/rand"

	"github.com/tjfoc/gmsm/sm4"
)

// PKCS7Padding 填充数据，使其符合块大小要求
func _PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 去除填充
func _PKCS7UnPadding(data []byte) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}

// 生成密钥
func GenerateKey() []byte {
	key := make([]byte, 16)
	rand.Read(key) // 从随机源中读取 16 字节
	return key
}

func Encrypt(content []byte, key []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 填充明文数据使其符合 SM4 的块大小要求（16 字节）
	content = _PKCS7Padding(content, block.BlockSize())

	// 加密
	data := make([]byte, len(content))
	block.Encrypt(data, content)

	return data, nil
}

func Decrypt(content []byte, key []byte) ([]byte, error) {
	// 创建一个 SM4 解密实例
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 解密
	data := make([]byte, len(content))
	block.Decrypt(data, content)

	// 去除填充
	data = _PKCS7UnPadding(data)

	return data, nil
}
