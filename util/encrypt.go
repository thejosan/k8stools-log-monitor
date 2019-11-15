package util

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/des"
	"errors"
	"bytes"
)

// MD5加密===================================================
func Md5_X_2(text string) string {
	return Md5(Md5(text))
}

func Md5(text string) string {
	str := md5.New()
	str.Write([]byte(text))
	return hex.EncodeToString(str.Sum(nil))
}

// 对称加密函数================================================
var keytxt string = "999898" //加密、解密秘钥 填8个字符
func Encrypt(text string) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher([]byte(keytxt))
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}
// 解密函数
func Decrypt(decrypted string) (string, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher([]byte(keytxt))
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	return string(out), nil
}
//填充
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}
//去填充
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
