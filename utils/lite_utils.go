package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"math/rand"
	"time"
)

// RandStr 指定长度随机字符串
func RandStr(length int, baseStrings ...string) string {
	var baseString string
	if len(baseStrings) > 0 && len(baseStrings[0]) > 0 {
		baseString = baseStrings[0]
	} else {
		baseString = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
	}
	b := make([]byte, length)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = baseString[rand.Intn(len(baseString))]
	}
	return string(b)
}

// AesCbcEncrypt AES-128-CBC-PKCS5Padding 加密主体
func AesCbcEncrypt(plainText []byte, key []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = Pkcs5Padding(plainText, 16)
	//指定初始向量vi
	//指定分组模式，返回一个BlockMode接口对象
	iv := []byte(RandStr(16))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

func AesCbcEncrypt7(plainText []byte, key []byte, iv []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = Pkcs7Padding(plainText)
	//指定初始向量vi
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

// Pkcs5Padding 用于AES padding 填充
func Pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func Pkcs7Padding(cipherText []byte) []byte {
	padding := 16 - len(cipherText)%16
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func EncryptCpdailyExtension(text []byte) string {
	desKey := "XCE927=="
	iv := "\x01\x02\x03\x04\x05\x06\x07\x08"
	res, err := DesEncryption([]byte(desKey), []byte(iv), text)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(res)
}

func EncryptBodyString(text []byte) string {
	aesKey := "SASEoK4Pa5d4SssO"
	iv := "\x01\x02\x03\x04\x05\x06\x07\x08\t\x01\x02\x03\x04\x05\x06\x07"
	res := AesCbcEncrypt7(text, []byte(aesKey), []byte(iv))
	return base64.StdEncoding.EncodeToString(res)
}

func DesEncryption(key, iv, plainText []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData := Pkcs5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}
