package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

func DecodePublicKey(encoded string) (publicKey *rsa.PublicKey, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	publicKey, err = x509.ParsePKCS1PublicKey(decoded)
	return
}

func EncodePublicKey(publicKey *rsa.PublicKey) string {
	marshalled := x509.MarshalPKCS1PublicKey(publicKey)
	encoded := base64.StdEncoding.EncodeToString(marshalled)
	return encoded
}

func CheckSignature(publicKey *rsa.PublicKey, signature []byte, payload []byte) error {
	hashed := sha256.Sum256(payload)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, nil)
}

func Sign(privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
	hashed := sha256.Sum256(payload)
	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], nil)
}
