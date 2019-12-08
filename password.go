package gorest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_#~()|[]?!*$^%:/;.,&{}=<>"

// GenerateSalt Generate a salt string with n length
func GenerateSalt(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// EncodePassword Encode password with the following algorithm:
// sha512(hmac-sha256(username + ":" + sha512(salt1) + sha512("/" + salt2) , password))
func EncodePassword(username string, password string, salt1 string, salt2 string) string {
	hasher := sha512.New()
	hasher.Write([]byte(salt1))

	var buffer bytes.Buffer
	buffer.WriteString(username)
	buffer.WriteString(":")
	buffer.Write(hasher.Sum(nil))
	hasher.Reset()

	hasher.Write([]byte("/" + salt2))
	buffer.Write(hasher.Sum(nil))
	hasher.Reset()

	mac := hmac.New(sha256.New, []byte(password))
	mac.Write(buffer.Bytes())

	// Hash hmac with sha512 algorithm
	hasher.Write(mac.Sum(nil))
	return hex.EncodeToString(hasher.Sum(nil))
}
