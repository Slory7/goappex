package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	padData := PKCS7Padding(origData, blockSize)
	crypted := make([]byte, blockSize+len(padData))
	iv := crypted[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted[blockSize:], padData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	iv := crypted[:blockSize]
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted)-blockSize)
	blockMode.CryptBlocks(origData, crypted[blockSize:])
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func AesEncryptBase64(origData, key string) (string, error) {
	data, err := AesEncrypt([]byte(origData), []byte(key))
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(data)
	return s, nil
}

func AesDecryptBase64(crypted, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}
	decrypted, err := AesDecrypt(data, []byte(key))
	if err != nil {
		return "", err
	}
	s := string(decrypted)
	return s, nil
}
