package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pass-it/internal/crypto"
	"strings"

	"github.com/google/uuid"
)

func main() {

	var hostAddr string
	var secret string
	flag.StringVar(&hostAddr, "host", "http://localhost:8080", "Pass It server address")
	flag.Parse()

	secret = flag.Arg(0)

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey

	encoded := crypto.EncodePublicKey(&publicKey)

	fmt.Println("Sending secret to server...")
	var id string = uuid.NewString()

	var storeParams = url.Values{"key": {encoded}, "payload": {secret}}

	var storeParamsEncoded = storeParams.Encode()

	storeRequest, err := http.NewRequest(http.MethodPut, hostAddr+"/store/"+id, strings.NewReader(storeParamsEncoded))
	if err != nil {
		panic(err)
	}
	storeRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(storeRequest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("Store request failed: %s/%d\n%s", resp.Status, resp.StatusCode, body))
	}

	fmt.Println("Sending secret to server... Done.")

	fmt.Println("Enter session id:")
	var sessionId string
	fmt.Scanf("%s", &sessionId)

	fmt.Println("Confirming secret receiver...")
	signature, err := crypto.Sign(privateKey, []byte(sessionId))
	if err != nil {
		panic(err)
	}
	signatureEncoded := base64.StdEncoding.EncodeToString(signature)

	var confirmParams = url.Values{"signature": {signatureEncoded}, "session_id": {sessionId}}

	var confirmParamsEncoded = confirmParams.Encode()

	confirmRequest, err := http.NewRequest(http.MethodPost, hostAddr+"/store/"+id, strings.NewReader(confirmParamsEncoded))
	if err != nil {
		panic(err)
	}
	confirmRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = http.DefaultClient.Do(confirmRequest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("Confirm request failed: %s/%d\n%s", resp.Status, resp.StatusCode, body))
	}
	fmt.Println("Confirming secret receiver... Done.")
}
