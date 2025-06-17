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

// MD5 returns the MD5 hash of the input string in hexadecimal format
func MD5(input string) string {
	return Md5(input)
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

// CustomEncrypt implements a custom encryption algorithm similar to the PHP version
func CustomEncrypt(text string, skey string) string {
	if skey == "" {
		skey = "key"
	}

	// Base64 encode the input string
	encoded := base64.StdEncoding.EncodeToString([]byte(text))

	// Convert to rune slice for easier manipulation
	strArr := []rune(encoded)
	strCount := len(strArr)

	// Append key characters to the base64 characters
	for i, v := range skey {
		if i < strCount {
			strArr[i] = strArr[i]
			strArr[i] = rune(strArr[i]) + v
		}
	}

	// Replace base64 characters with custom strings
	result := string(strArr)
	result = strings.ReplaceAll(result, "=", "O0O0O")
	result = strings.ReplaceAll(result, "+", "o000o")
	result = strings.ReplaceAll(result, "/", "oo00o")

	return result
}

// CustomDecrypt implements a custom decryption algorithm similar to the PHP version
func CustomDecrypt(text string, skey string) string {
	if skey == "" {
		skey = "key"
	}

	// Replace custom strings back to base64 characters
	result := text
	result = strings.ReplaceAll(result, "O0O0O", "=")
	result = strings.ReplaceAll(result, "o000o", "+")
	result = strings.ReplaceAll(result, "oo00o", "/")

	// Split into pairs of characters
	var strArr []string
	for i := 0; i < len(result); i += 2 {
		if i+1 < len(result) {
			strArr = append(strArr, result[i:i+2])
		}
	}

	// Remove key characters
	for i, v := range skey {
		if i < len(strArr) {
			if len(strArr[i]) > 1 && strArr[i][1] == byte(v) {
				strArr[i] = string(strArr[i][0])
			}
		}
	}

	// Join and decode
	decoded, err := base64.StdEncoding.DecodeString(strings.Join(strArr, ""))
	if err != nil {
		return ""
	}

	return string(decoded)
}
