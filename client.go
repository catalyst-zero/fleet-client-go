package client

import (
	"github.com/juju/errgo"

	execPkg "os/exec"
)

const (
	FLEETCTL        = "fleetctl"
	ENDPOINT_OPTION = "--endpoint"
	ENDPOINT_VALUE  = "http://172.17.42.1:4001"
)

type Client struct {
	etcdPeer string
}

func NewClient() *Client {
	return &Client{
		etcdPeer: ENDPOINT_VALUE,
	}
}

func (this *Client) SetEtcdPeer(etcdPeer string) {
	this.etcdPeer = etcdPeer
}

func (this *Client) Submit(filePath string) error {
	cmd := execPkg.Command(FLEETCTL, ENDPOINT_OPTION, this.etcdPeer, "submit", filePath)
	_, err := exec(cmd)

	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (this *Client) Start(unitFileName string) error {
	cmd := execPkg.Command(FLEETCTL, ENDPOINT_OPTION, this.etcdPeer, "start", "--no-block=true", unitFileName)
	_, err := exec(cmd)

	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (this *Client) Stop(unitFileName string) error {
	cmd := execPkg.Command(FLEETCTL, ENDPOINT_OPTION, this.etcdPeer, "stop", "--no-block=true", unitFileName)
	_, err := exec(cmd)

	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (this *Client) Destroy(unitFileName string) error {
	cmd := execPkg.Command(FLEETCTL, ENDPOINT_OPTION, this.etcdPeer, "destroy", unitFileName)
	_, err := exec(cmd)

	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (this *Client) Status(unitFileName string) (Status, error) {
	cmd := execPkg.Command(FLEETCTL, ENDPOINT_OPTION, this.etcdPeer, "list-units")
	stdout, err := exec(cmd)

	if err != nil {
		return Status{}, errgo.Mask(err)
	}

	running, err := isRunning(unitFileName, stdout)
	if err != nil {
		return Status{}, errgo.Mask(err)
	}

	return Status{Running: running, ContainerIP: "127.0.0.1"}, nil
}
