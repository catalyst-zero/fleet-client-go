package main

import (
	"fmt"
	"os"

	"github.com/catalyst-zero/fleet-client-go"
)

func main() {
	// cAPI := client.NewClientCLIWithPeer("http://127.0.0.1:4001")
	cAPI := client.NewClientAPI()

	j, err := cAPI.StatusUnit(os.Args[1])
	fmt.Println(j, err)
}
