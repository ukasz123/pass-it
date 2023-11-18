package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"pass-it/crypto"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey

	encoded := crypto.EncodePublicKey(&publicKey)
	
	fmt.Println("Encoded public key:", encoded)

	var sessionId string
	fmt.Scanf("%s", &sessionId)

	signature, err := crypto.Sign(privateKey, []byte(sessionId))
	if err != nil {
		panic(err)
	}
	signatureEncoded := base64.StdEncoding.EncodeToString(signature)
	fmt.Println("Signature:", signatureEncoded)
}
