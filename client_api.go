package client

import (
	"net"
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/unit"
	"github.com/coreos/fleet/job"
	"github.com/juju/errgo"
)

const (
	FLEET_SOCK = "/var/run/fleet.sock"
)

type ClientAPI struct {
	client client.API
}

func NewClientAPI() FleetClient {
	return NewClientAPIWithSocket(FLEET_SOCK)
}

func NewClientAPIWithSocket(socket string) FleetClient {
	dialFunc := func(string, string) (net.Conn, error) {
		return net.Dial("unix", socket)
	}

	trans := http.Transport{
		Dial: dialFunc,
	}

	hc := http.Client{
		Transport: &trans,
	}
	
	c, _ := client.NewHTTPClient(&hc)

	return &ClientAPI{
		client: c,
	}
}

// getUnitFromFile attempts to load a Unit from a given filename
// It returns the Unit or nil, and any error encountered
func getUnitFromFile(file string) (*unit.Unit, error) {
	out, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return unit.NewUnit(string(out))
}

func (this *ClientAPI) Submit(name, filePath string) error {
	unit, err := getUnitFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed reading unit file %s: %v", name, err)
	}
	
	j := job.NewJob(name, *unit)

	if err := this.client.CreateJob(j); err != nil {
		return fmt.Errorf("failed creating job %s: %v", j.Name, err)
	}

	return nil
}

func (this *ClientAPI) Get(name string) (*job.Job, error) {
	j, err := this.client.Job(name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Job(%s) from Registry: %v", name, err)
	} else if j == nil {
		return nil, fmt.Errorf("unable to find Job(%s)", name)
	}
	return j, nil
}

func (this *ClientAPI) Start(name string) error {
	j, err := this.Get(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetJobTargetState(j.Name, job.JobStateLaunched)
	
	return nil
}

func (this *ClientAPI) Stop(name string) error {
	j, err := this.Get(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetJobTargetState(j.Name, job.JobStateLoaded)
	
	return nil
}

func (this *ClientAPI) Destroy(name string) error {
	j, err := this.Get(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.DestroyJob(j.Name)
	
	return nil
}

func (this *ClientAPI) Status(name string) (*Status, error) {
	return nil, fmt.Errorf("Method not implemented: ClientAPI.Status")
}
