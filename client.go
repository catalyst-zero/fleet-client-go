package client

import (
	"github.com/coreos/fleet/job"
)

type Status struct {
	Running     bool
	ContainerIP string
}

type FleetClient interface {
	Get(name string) (*job.Job, error)
	Submit(name, filePath string) error
	Start(name string) error
	Stop(name string) error
	Destroy(name string) error
	Status(name string) (Status, error)
}

func NewClient() FleetClient {
	return NewClientCLI()
}
