package client

import (
	"bytes"
	"os/exec"
)

const (
	FLEETCTL        = "fleetctl"
	ENDPOINT_OPTION = "--endpoint"
	ENDPOINT_VALUE  = "http://172.17.42.1:4001"
)

type Status struct {
	Running     bool
	ContainerIP string
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (this *Client) Submit(filePath string) error {
	cmd := exec.Command(FLEETCTL, ENDPOINT_OPTION, ENDPOINT_VALUE, "submit", filePath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (this *Client) Start(unitFileName string) error {
	cmd := exec.Command(FLEETCTL, ENDPOINT_OPTION, ENDPOINT_VALUE, "start", unitFileName)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (this *Client) Stop(unitFileName string) error {
	cmd := exec.Command(FLEETCTL, ENDPOINT_OPTION, ENDPOINT_VALUE, "stop", unitFileName)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (this *Client) Destroy(unitFileName string) error {
	cmd := exec.Command(FLEETCTL, ENDPOINT_OPTION, ENDPOINT_VALUE, "destroy", unitFileName)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (this *Client) Status(unitFileName string) (Status, error) {
	var out bytes.Buffer

	cmd := exec.Command(FLEETCTL, ENDPOINT_OPTION, ENDPOINT_VALUE, "list-units")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return Status{}, err
	}

	running, err := isRunning(unitFileName, out.String())
	if err != nil {
		return Status{}, err
	}

	return Status{Running: running, ContainerIP: "127.0.0.1"}, nil
}
