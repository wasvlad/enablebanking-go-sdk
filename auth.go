package enablebankinggosdk

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyContent)
	if block == nil {
		return nil, errors.New("No PEM data found in file")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return privateKey, nil
	case "PRIVATE KEY":
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return privateKey.(*rsa.PrivateKey), nil
	default:
		return nil, errors.New("Unsupported type of private key")
	}
}

type JWTSettings struct{
	key *rsa.PrivateKey
	issuer string
	audience string
	valid_for int64
	kid string
}

func generateJWTToken(settings JWTSettings) (string, error) {
	iat := time.Now().UTC().Unix()
	claims := jwt.MapClaims{
		"iss": settings.issuer,
		"aud": settings.audience,
		"iat": iat,
		"exp": iat + settings.valid_for,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if token.Header == nil {
		token.Header = make(map[string]interface{})
	}

	token.Header["kid"] = settings.kid
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "RS256"

	signed_token, err := token.SignedString(settings.key)

	if err != nil {
		return "", err
	}

	return signed_token, nil
}
