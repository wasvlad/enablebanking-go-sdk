package eb

type EBClient struct {
	jwtToken     string
	JWTGenerator jwtGenerator
	httpClient   httpClient
}

const apiURL string = "https://api.enablebanking.com"

func NewEBClient(jwtGenerator jwtGenerator) (*EBClient, error) {
	jwtToken, err := jwtGenerator.generate()
	if err != nil {
		return nil, err
	}

	ebClient := EBClient{
		jwtToken:     jwtToken,
		JWTGenerator: jwtGenerator,
		httpClient:   newHttpClient(),
	}

	return &ebClient, nil
}
