// Copyright 2019 py60800.
// Use of this source code is governed by Apache-2 licence

package tuya

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func md5Sign(b []byte, key []byte, version string) []byte {
	h := md5.New()
	h.Write([]byte("data="))
	h.Write(b)
	h.Write([]byte("||lpv=" + version + "||"))
	h.Write(key)

	hash := h.Sum(nil)
	// 8:16
	return []byte(hex.EncodeToString(hash[4:12]))
}

func aesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	remain := len(data) % blockSize
	if remain == 0 {
		remain = blockSize
	}

	padding := make([]byte, blockSize-remain)
	for i := range padding {
		padding[i] = byte(blockSize - remain)
	}
	data = append(data, padding...)

	ciphertext := make([]byte, len(data))
	for i := 0; i < len(data); i = i + blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], data[i:i+blockSize])
	}

	//return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
	return ciphertext, nil
}

func base64decode(data []byte) []byte {
	dstLen := base64.StdEncoding.DecodedLen(len(data))
	dst := make([]byte, dstLen)
	_, err := base64.StdEncoding.Decode(dst, data)
	if err != nil {
		return nil
	}
	return dst
}

func aesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	length := len(ciphertext)
	if length < 16 {
		return nil, errors.New("encrypted block is too small")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	blockSize := block.BlockSize()
	if length%blockSize != 0 {
		return []byte{}, errors.New("Bad ciphertext len")
	}

	cleartext := make([]byte, length)
	for i := 0; i < length; i = i + blockSize {
		block.Decrypt(cleartext[i:i+blockSize], ciphertext[i:i+blockSize])
	}

	// remove padding
	padding := int(cleartext[length-1])
	if padding < 0 || padding > blockSize {
		return []byte{}, errors.New("Bad padding")
	}

	return cleartext[:length-padding], nil
}
