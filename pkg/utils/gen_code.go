package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateCaptcha(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	alphabets := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	captcha := ""
	for i := 0; i < length; i++ {
		switch rand.Intn(2) {
		case 0:
			captcha += string(digits[rand.Intn(len(digits))])
		case 1:
			captcha += string(alphabets[rand.Intn(len(alphabets))])
		}

	}

	return captcha
}

func GenerateCode(length int) string {
	numeric := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// GenerateUserCode generates a unique code based on username and timestamp
// The format is like: 835149712G3A2G7C63897A168M7
func GenerateUserCode(username string, timestamp int64) string {
	// Use timestamp as seed to make it deterministic for same inputs
	rand.Seed(timestamp)

	// Define the character sets
	digits := "0123456789"
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var sb strings.Builder

	// Add 9 random digits
	for i := 0; i < 9; i++ {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}

	// Add 15 characters alternating between letters and digits
	for i := 0; i < 15; i++ {
		if i%2 == 0 {
			sb.WriteByte(letters[rand.Intn(len(letters))])
		} else {
			sb.WriteByte(digits[rand.Intn(len(digits))])
		}
	}

	return sb.String()
}
