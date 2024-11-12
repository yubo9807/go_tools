package m2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/tjfoc/gmsm/sm2"
)

// 解析私钥
func parsePrivateKey(privStr string) (*sm2.PrivateKey, error) {
	// 解码 Base64 字符串为字节数组
	privBytes, err := hex.DecodeString(privStr)
	if err != nil {
		return nil, err
	}

	// 使用解码后的字节数组构造私钥
	priv := new(sm2.PrivateKey)
	priv.D = new(big.Int).SetBytes(privBytes)
	priv.PublicKey = sm2.PublicKey{
		Curve: sm2.P256Sm2(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}

	return priv, nil
}

// 解析公钥
func parsePublicKey(pubStr string) (*sm2.PublicKey, error) {
	// 解码 Base64 字符串为字节数组
	pubBytes, err := hex.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}

	// 使用解码后的字节数组构造公钥
	pub := new(sm2.PublicKey)
	pub.Curve = sm2.P256Sm2()
	pub.X = new(big.Int).SetBytes(pubBytes[:32])
	pub.Y = new(big.Int).SetBytes(pubBytes[32:])
	return pub, nil
}

// 生成公私钥对
func GenerateKeyPair() (*string, *string, error) {
	// 生成私钥
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// 导出公钥
	public := priv.Public().(*sm2.PublicKey)

	b := append(public.X.Bytes(), public.Y.Bytes()...)

	privKey := hex.EncodeToString(priv.D.Bytes())
	pubKey := hex.EncodeToString(b)

	return &privKey, &pubKey, nil
}

// SM2 加密文件
func EncryptFile(filename, newFilename, publicStr string) error {

	// 读取原始文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取文件失败", err)
		return err
	}

	// 解析公钥
	publicKey, err := parsePublicKey(publicStr)
	if err != nil {
		fmt.Println("解析公钥失败", err)
		return err
	}

	// 使用 SM2 公钥加密内容
	newContent, err := sm2.Encrypt(publicKey, content, rand.Reader, 0)
	if err != nil {
		fmt.Println("加密失败", err)
		return err
	}

	// 将加密后的内容写入新文件
	if err := os.WriteFile(newFilename, newContent, 0644); err != nil {
		fmt.Println("写入文件失败", err)
		return err
	}

	return nil
}

// SM2 解密文件
func DecryptFile(filename, newFilename string, privateStr string) error {

	// 读取加密文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取文件失败", err)
		return err
	}

	// 解析 SM2 私钥
	privateKey, err := parsePrivateKey(privateStr)
	if err != nil {
		fmt.Println("解析私钥失败", err)
		return err
	}

	// 使用 SM2 私钥解密内容
	newContent, err := sm2.Decrypt(privateKey, content, 0)
	if err != nil {
		fmt.Println("解密失败", err)
		return err
	}

	// 将解密后的内容写入新文件
	if err := os.WriteFile(newFilename, newContent, 0644); err != nil {
		fmt.Println("写入文件失败", err)
		return err
	}

	return nil
}
