package m2

import (
	"crypto/rand"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
)

var mode = sm2.C1C3C2

// 生成密钥对
func GenerateKey() ([]byte, []byte, error) {
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// 导出公钥
	public := priv.Public().(*sm2.PublicKey)

	publicBytes := append(public.X.Bytes(), public.Y.Bytes()...)

	return priv.D.Bytes(), publicBytes, nil
}

// 加密
func Encrypt(content []byte, publicKey []byte) ([]byte, error) {
	pub := new(sm2.PublicKey)
	pub.Curve = sm2.P256Sm2()
	pub.X = new(big.Int).SetBytes(publicKey[:32])
	pub.Y = new(big.Int).SetBytes(publicKey[32:])

	return sm2.Encrypt(pub, content, rand.Reader, mode)
}

// 解密
func Decrypt(content []byte, privateKey []byte) ([]byte, error) {
	priv := new(sm2.PrivateKey)
	priv.D = new(big.Int).SetBytes(privateKey)
	priv.PublicKey = sm2.PublicKey{
		Curve: sm2.P256Sm2(),
		X:     new(big.Int),
		Y:     new(big.Int),
	}

	return sm2.Decrypt(priv, content, mode)
}
