package main

import (
	"fmt"
	"os"
	"time"

	eb "github.com/wasvlad/enablebanking-go-sdk"
)

func main() {
	if len(os.Args) < 4 {
		panic("Usage: " + os.Args[0] + " path_to_private_key issuer kid")
	}
	generator, err := eb.NewJWTSession(
		os.Args[1], // path
		os.Args[2], // issuer
		60*60,      // 1 hour
		os.Args[3], // kid
	)
	if err != nil {
		panic(err)
	}
	client, err := eb.NewEBClient(
		generator,
	)
	if err != nil {
		panic(err)
	}
	access := eb.NewAccess(time.Hour)
	aspsp := eb.Aspsp{
		Name:    "Bank Millennium",
		Country: "PL",
	}
	data := eb.StartAuthorizationRequest{
		Access:      access,
		Aspsp:       aspsp,
		RedirectUrl: "http://localhost",
		State:       "state",
	}
	response, err := client.Authorize(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Authorize here and provide the code " + response.URL)
}
