package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"log"
)

// Encrypt plain text using public key
func encryptText(plain_text string, public_key *rsa.PublicKey) string {

	cipher_text, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, public_key, []byte(plain_text), []byte(""))

	if err != nil {
		log.Fatalln("Error during plain text encryption: ", err.Error)
	}

	return string(cipher_text)
}

// Verify authentification code using public key
func verifAuthCode(authCode string, public_key *rsa.PublicKey, signature []byte) bool {

	hash := crypto.SHA512
	pssh := hash.New()
	pssh.Write([]byte(authCode))
	hashed := pssh.Sum(nil)

	err := rsa.VerifyPSS(public_key, hash, hashed, signature, &rsa.PSSOptions{})

	if err != nil {
		return false
	} else {
		return true
	}
}
