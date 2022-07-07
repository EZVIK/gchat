package util

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
)

//func GetPrivateKey(key string) (*rsa.PrivateKey, []byte, error) {
//	decoded, err := base64.StdEncoding.DecodeString(key)
//	if err != nil {
//		return nil, nil, err
//	}
//	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(decoded)
//	if err != nil {
//		return nil, nil, err
//	}
//	return priKey, decoded, err
//}

func GetPrivateKey(key string) (*rsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(decoded)
	return priKey, err
}
