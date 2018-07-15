package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"os"
)

// Encrypt plain text using public key
func encryptText(plainTextMsg string, clientPublicKey *rsa.PublicKey) string {

	cipherTextMsg, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, clientPublicKey, []byte(plainTextMsg), []byte(""))

	if err != nil {
		fmt.Println("Error during plain text encryption: ", err.Error)
		os.Exit(1)
	}

	return string(cipherTextMsg)
}

// Verify authentification code using public key
func verifAuthCode(authCode string, clientPublicKey *rsa.PublicKey, signature []byte) bool {

	hash := crypto.SHA512
	pssh := hash.New()
	pssh.Write([]byte(authCode))
	hashed := pssh.Sum(nil)

	err := rsa.VerifyPSS(clientPublicKey, hash, hashed, signature, &rsa.PSSOptions{})

	if err != nil {
		return false
	} else {
		return true
	}
}
