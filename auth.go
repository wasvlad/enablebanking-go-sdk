package eb

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtGenerator struct {
	key      *rsa.PrivateKey
	issuer   string
	validFor int64
	kid      string
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key from disk: %v", err)
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

func NewJWTGenerator(path string, issuer string, validFor int64, kid string) (*jwtGenerator, error) {
	key, err := loadPrivateKey(path)
	if err != nil {
		return nil, err
	}
	return &jwtGenerator{
		key:      key,
		issuer:   issuer,
		validFor: validFor,
		kid:      kid,
	}, nil
}

func (generator jwtGenerator) generate() (string, error) {
	iat := time.Now().UTC().Unix()
	claims := jwt.MapClaims{
		"iss": generator.issuer,
		"aud": strings.Split(apiURL, "//")[1],
		"iat": iat,
		"exp": iat + generator.validFor,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if token.Header == nil {
		token.Header = make(map[string]interface{})
	}

	token.Header["kid"] = generator.kid
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "RS256"

	signed_token, err := token.SignedString(generator.key)

	if err != nil {
		return "", err
	}

	return signed_token, nil
}

type access struct {
	/*
		This parameter specifies the date and time until which the authorised session
		remains valid. The value must be in the RFC3339 date-time format with a timezone
		offset, e.g. 2025-12-01T12:00:00.000000+00:00. The provided value cannot exceed
		the date and time, calculated as "now" + maximum_consent_validity (provided in
		seconds for each ASPSP in response to the GET /aspsps request). The provided value
		is subject to adjustment to comply with the ASPSP's requirements. Specifically, if
		the provided value is less than the minimum consent validity allowed by the ASPSP
		(e.g., some ASPSPs require a minimum of 1 hour or 1 day), the consent validity will
		be adjusted to meet these requirements. However, the session validity will remain
		exactly as specified. This means that even if the consent remains valid on the
		ASPSP's side, the session will expire based on the initially provided value.
	*/
	ValidUntil string `json:"valid_until"`
}

func NewAccess(validFor time.Duration) access {
	validUntil := time.Now().UTC().Add(validFor).Format(time.RFC3339)
	return access{
		ValidUntil: validUntil,
	}
}

type Aspsp struct {
	// Name of the ASPSP (i.e. a bank or a similar financial institution)
	Name string `json:"name"`
	// Two-letter ISO 3166 code of the country, in which ASPSP operates
	Country string `json:"country"`
}

type StartAuthorizationRequest struct {
	// Scope of access to be request from ASPSP and to be confirmed by PSU
	Access access `json:"access"`
	// ASPSP that PSU is going to be authenticated to
	Aspsp Aspsp `json:"aspsp"`
	// Arbitrary string. Same string will be returned in query parameter when redirecting to the URL passed via redirect_url parameter
	State string `json:"state"`
	// URL that PSU will be redirected to after authorization
	RedirectUrl string `json:"redirect_url"`
}

func (ebClient EBClient) Authorize(data StartAuthorizationRequest) (*http.Response, error) {
	jsonBytes, err := json.Marshal(data)
	buf := bytes.NewBuffer(jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}
	request, err := http.NewRequest(
		"POST",
		apiURL+"/auth",
		buf,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+ebClient.jwtToken)
	request.Header.Set("Content-Type", "application/json")
	resp, err := ebClient.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed executing request: %v", err)
	}

	return resp, nil
}
