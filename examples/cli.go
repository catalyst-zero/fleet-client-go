package main

import (
	"fmt"
	"github.com/catalyst-zero/fleet-client-go"
)

func main() {
	cAPI := client.NewClientAPI()

	j, err := cAPI.Get("app.service")
	fmt.Println(j, err)
}
