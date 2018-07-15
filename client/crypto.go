package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"os"
)

// Generate RSA-2048 key pair
func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {

	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println("Error during key pair generation: ", err.Error)
		os.Exit(1)
	} else {
		fmt.Println("Key pair generetad successfully")
	}

	return clientPrivateKey, &clientPrivateKey.PublicKey
}

// Decrypt ciphertext using private key
func decryptText(cipherTextMsg string, clientPrivateKey *rsa.PrivateKey) string {

	decryptedTextMsg, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, clientPrivateKey, []byte(cipherTextMsg), []byte(""))

	if err != nil {
		fmt.Println("Error during plain text decryption: ", err.Error)
		os.Exit(1)
	}

	return string(decryptedTextMsg)
}

// Sign authentification code using private key
func signAuthCode(authCode string, clientPrivateKey *rsa.PrivateKey) []byte {

	hash := crypto.SHA512
	pssh := hash.New()
	pssh.Write([]byte(authCode))
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, clientPrivateKey, hash, hashed, &rsa.PSSOptions{})

	if err != nil {
		fmt.Println("Error during authentification code signing: ", err.Error)
		os.Exit(1)
	}

	return signature
}
