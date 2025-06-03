package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
)

var tBytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Md5(password string) string {
	pw := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return pw
}

func GenerateSalt(mid, userName string) string {
	h := md5.New()
	io.WriteString(h, mid+userName)
	salt := fmt.Sprintf("%x", h.Sum(nil))[:16]
	return strings.ToUpper(salt)
}

func HashPassword(salt, password string) string {
	h := md5.New()
	io.WriteString(h, salt+password)
	pw := fmt.Sprintf("%x", h.Sum(nil))
	return strings.ToUpper(pw)
}

func GenAccessToken() string {
	token, _ := uuid.NewRandom()
	return Encode([]byte(token.String()))
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, tBytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, tBytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func HashPwd(username, password string) (string, string) {
	salt := GenerateSalt("588326785867908888", username)
	return salt, HashPassword(salt, password)
}
