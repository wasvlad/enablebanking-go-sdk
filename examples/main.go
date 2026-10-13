package main

import (
	"fmt"
	"io"
	"os"
	"time"

	ebSDK "github.com/wasvlad/enablebanking-go-sdk"
)

func main() {
	if len(os.Args) < 4 {
		panic("Usage: " + os.Args[0] + " path_to_private_key issuer kid")
	}
	generator, err := ebSDK.NewJWTGenerator(
		os.Args[1],
		os.Args[2],
		60*60,
		os.Args[3],
	)
	if err != nil {
		panic(err)
	}
	client, err := ebSDK.NewEBClient(
		*generator,
	)
	if err != nil {
		panic(err)
	}
	access := ebSDK.NewAccess(time.Hour)
	aspsp := ebSDK.Aspsp{
		Name:    "Bank Millennium",
		Country: "PL",
	}
	data := ebSDK.StartAuthorizationRequest{
		Access:      access,
		Aspsp:       aspsp,
		RedirectUrl: "http://localhost",
		State:       "state",
	}
	response, err := client.Authorize(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("end")
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body[:]))
}
