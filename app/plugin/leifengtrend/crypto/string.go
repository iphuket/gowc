package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// DeSting string use aes decrypt
func DeSting(str, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(str))
	blockMode.CryptBlocks(origData, str)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// EnSting string use aes encrypt
func EnSting(str, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	str = PKCS5Padding(str, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	enStr := make([]byte, len(str))
	blockMode.CryptBlocks(enStr, str)
	return enStr, nil
}

// PKCS5Padding func
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding func
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
