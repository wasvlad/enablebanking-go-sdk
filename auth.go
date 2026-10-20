package eb

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtToken struct {
	token      string
	validUntil int64
}

func (token *jwtToken) isValid(tolerance int64) bool {
	iat := time.Now().Unix()
	if token.validUntil-iat < tolerance {
		return false
	}
	return true
}

type jwtGenerator struct {
	key      *rsa.PrivateKey
	issuer   string
	validFor int64
	kid      string
	token    jwtToken
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key from disk: %v", err)
	}
	block, _ := pem.Decode(keyContent)
	if block == nil {
		return nil, errors.New("no PEM data found in file")
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
		return nil, errors.New("unsupported type of private key")
	}
}

func NewJWTGenerator(path string, issuer string, validFor int64, kid string) (*jwtGenerator, error) {
	key, err := loadPrivateKey(path)
	if err != nil {
		return nil, err
	}
	token := jwtToken{
		token:      "",
		validUntil: time.Now().Unix() - 1,
	}
	return &jwtGenerator{
		key:      key,
		issuer:   issuer,
		validFor: validFor,
		kid:      kid,
		token:    token,
	}, nil
}

func (generator *jwtGenerator) generate() error {
	iat := time.Now().UTC().Unix()
	exp := iat + generator.validFor
	claims := jwt.MapClaims{
		"iss": generator.issuer,
		"aud": strings.Split(apiURL, "//")[1],
		"iat": iat,
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if token.Header == nil {
		token.Header = make(map[string]interface{})
	}

	token.Header["kid"] = generator.kid
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "RS256"

	signedToken, err := token.SignedString(generator.key)

	if err != nil {
		return err
	}

	generator.token.token = signedToken
	generator.token.validUntil = exp

	return nil
}

func (gen *jwtGenerator) getToken() (string, error) {
	var err error = nil
	if !gen.token.isValid(10) {
		err = gen.generate()
	}
	if err != nil {
		return "", fmt.Errorf("unable to generate token: %v", err)
	}
	return gen.token.token, nil
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
	token, err := ebClient.JWTGenerator.getToken()
	if err != nil {
		return nil, err
	}
	r := ebClient.httpClient.R().bearerAuthorization(token)
	r, err = r.jsonBody(data)
	if err != nil {
		return nil, fmt.Errorf("error creating json body: %v", err)
	}
	resp, err := r.post(apiURL + "/auth")
	if err != nil {
		return nil, fmt.Errorf("error making requst: %v", err)
	}
	return resp, nil
}
