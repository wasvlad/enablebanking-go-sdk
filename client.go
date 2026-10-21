package eb

import "fmt"

type EBClient struct {
	jwtSession *jwtSession
	httpClient httpClient
}

const apiURL string = "https://api.enablebanking.com"

func NewEBClient(jwtSession *jwtSession) (*EBClient, error) {
	err := jwtSession.refresh()
	if err != nil {
		return nil, fmt.Errorf("unable to generate jwt token: %w", err)
	}

	ebClient := EBClient{
		jwtSession: jwtSession,
		httpClient: newHttpClient(),
	}

	return &ebClient, nil
}
