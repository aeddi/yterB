package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"log"
)

// Generate RSA-2048 key pair
func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {

	private_key, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatal("Error during key pair generation: ", err)
	} else {
		log.Println("Key pair generetad successfully")
	}

	return private_key, &private_key.PublicKey
}

// Decrypt ciphertext using private key
func decryptText(cipherTextMsg string, private_key *rsa.PrivateKey) string {

	decrypted_text, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, private_key, []byte(cipherTextMsg), []byte(""))

	if err != nil {
		log.Fatalln("Error during plain text decryption: ", err)
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
		log.Fatalln("Error during authentification code signing: ", err)
	}

	return signature
}
