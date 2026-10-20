package eb

import "fmt"

type EBClient struct {
	JWTGenerator *jwtGenerator
	httpClient   httpClient
}

const apiURL string = "https://api.enablebanking.com"

func NewEBClient(jwtGenerator *jwtGenerator) (*EBClient, error) {
	err := jwtGenerator.generate()
	if err != nil {
		return nil, fmt.Errorf("unable to generate jwt token: %v", err)
	}

	ebClient := EBClient{
		JWTGenerator: jwtGenerator,
		httpClient:   newHttpClient(),
	}

	return &ebClient, nil
}
