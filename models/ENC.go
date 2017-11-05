package models

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strconv"
)

func Test(str string) {
	aesEnc := AesEncrypt{"1234567890123456"}
	arrEncrypt, err := aesEnc.Encrypt(str)
	if err != nil {
		fmt.Println(arrEncrypt)
		return
	}
	fmt.Println("Enced:", arrEncrypt)
	strMsg, err := aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		fmt.Println(arrEncrypt)
		return
	}
	fmt.Println(strMsg)
}

var ENC = AesEncrypt{"0123456789012345"}

type AesEncrypt struct {
	strKey string
}

func (this *AesEncrypt) getKey() []byte {
	keyLen := len(this.strKey)
	if keyLen < 16 {
		this.strKey = "0123456789012345"
	}
	arrKey := []byte(this.strKey)
	if keyLen >= 32 {
		return arrKey[:32]
	}
	if keyLen >= 24 {
		return arrKey[:24]
	}
	return arrKey[:16]
}

func (this *AesEncrypt) Encrypt(strMesg string) (string, error) {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(strMesg))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strMesg))

	encstr := ""
	for _, v := range encrypted {
		encstr += fmt.Sprintf("%02x", v)
	}

	return encstr, nil
}

func (this *AesEncrypt) Decrypt(srcStr string) (strDesc string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src := []byte{}
	var b int64
	for i := 0; i < len(srcStr)-1; i += 2 {
		b, _ = strconv.ParseInt(srcStr[i:i+2], 16, 8)
		src = append(src, byte(b))
	}

	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}
