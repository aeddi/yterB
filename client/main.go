package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"os"
	"bufio"
	"strings"
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

// Get user entry from stdin
func getUserEntry(prompt string) string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

// Encrypt plain text using public key
func encryptText(plainTextMsg string, clientPublicKey *rsa.PublicKey) string {

	cipherTextMsg, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, clientPublicKey, []byte(plainTextMsg), []byte(""))

	if err != nil {
		fmt.Println("Error during plain text encryption: ", err.Error)
		os.Exit(1)
	}

	return string(cipherTextMsg)
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

func main() {

	clientPrivateKey, clientPublicKey := generateKeyPair()

	plainTextMsg := getUserEntry("Enter text to encrypt/decrypt: ")
	cipherTextMsg := encryptText(plainTextMsg, clientPublicKey)
	decryptedTextMsg := decryptText(cipherTextMsg, clientPrivateKey)

	if plainTextMsg == decryptedTextMsg {
		fmt.Println("Encryption/decryption works!")
		fmt.Println(plainTextMsg, "==", decryptedTextMsg)
	}

	authCode := getUserEntry("Enter text to sign/verify: ")
	signature := signAuthCode(authCode, clientPrivateKey)

	if verifAuthCode(authCode, clientPublicKey, signature) {
		fmt.Println("Siging/verifying works!")
		fmt.Println(authCode)
	}
}
