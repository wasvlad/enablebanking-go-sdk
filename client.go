package eb

import "net/http"

const apiURL string = "https://api.enablebanking.com"

type EBClient struct {
	jwtToken     string
	JWTGenerator jwtGenerator
	httpClient   *http.Client
}

func NewEBClient(jwtGenerator jwtGenerator) (*EBClient, error) {
	jwtToken, err := jwtGenerator.generate()
	if err != nil {
		return nil, err
	}

	ebClient := EBClient{
		jwtToken:     jwtToken,
		JWTGenerator: jwtGenerator,
		httpClient:   http.DefaultClient,
	}

	return &ebClient, nil
}
