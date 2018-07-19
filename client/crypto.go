package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"os"
)

// Generate RSA-2048 key pair
func generateKeyPair() *rsa.PrivateKey {

	key_pair, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		consoleLog("Error during key pair generation: " + err.Error())
		os.Exit(1)
	} else {
		consoleLog("Key pair generetad successfully")
	}

	return key_pair
}

// Decrypt ciphertext using private key
func decryptText(cipherTextMsg string, private_key *rsa.PrivateKey) string {

	decrypted_text, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, private_key, []byte(cipherTextMsg), []byte(""))

	if err != nil {
		consoleLog("Error during plain text decryption: " + err.Error())
		os.Exit(1)
	}

	return string(decrypted_text)
}

// Sign authentification code using private key
func signAuthCode(authCode string, private_key *rsa.PrivateKey) []byte {

	hash := crypto.SHA512
	pssh := hash.New()
	pssh.Write([]byte(authCode))
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, private_key, hash, hashed, &rsa.PSSOptions{})

	if err != nil {
		consoleLog("Error during authentification code signing: " + err.Error())
		os.Exit(1)
	}

	return signature
}
