package main

import (
	"fmt"
	"github.com/catalyst-zero/fleet-client-go"
)

func main() {
	cAPI := client.NewClientAPI()

	u, err := cAPI.Unit("app.service")
	fmt.Println(u, err)
}
