package models

import "crypto/rsa"

type StoredData[T any, P any] struct {
	Key T;
	Payload P;
}

type DefaultStoredData = StoredData[*rsa.PublicKey, string]


type PayloadMessage[Addr any, Payload any] struct {
	Addr Addr;
	Payload Payload;
}
