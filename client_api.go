package client

import (
	"fmt"
	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
	"github.com/juju/errgo"
	"io/ioutil"
	"net"
	"net/http"
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
func getUnitFromFile(file string) (*unit.UnitFile, error) {
	out, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return unit.NewUnitFile(string(out))
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

func (this *ClientAPI) ScheduledUnit(name string) (*job.ScheduledUnit, error) {
	su, err := this.client.ScheduledUnit(name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ScheduledUnit (%s) from Registry: %v", name, err)
	} else if su == nil {
		return nil, fmt.Errorf("unable to find ScheduledUnit (%s)", name)
	}
	return su, nil
}

func (this *ClientAPI) Unit(name string) (*job.Unit, error) {
	u, err := this.client.Unit(name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Unit (%s) from Registry: %v", name, err)
	} else if u == nil {
		return nil, fmt.Errorf("unable to find Unit (%s)", name)
	}
	return u, nil
}

func (this *ClientAPI) Start(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetJobTargetState(u.Name, job.JobStateLaunched)

	return nil
}

func (this *ClientAPI) Stop(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetJobTargetState(u.Name, job.JobStateLoaded)

	return nil
}

func (this *ClientAPI) Load(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetJobTargetState(u.Name, job.JobStateLoaded)

	return nil
}

func (this *ClientAPI) Destroy(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.DestroyJob(u.Name)

	return nil
}

func (this *ClientAPI) Status(name string) (*Status, error) {
	return nil, fmt.Errorf("Method not implemented: ClientAPI.Status")
}

func (this *ClientAPI) StatusUnit(name string) (UnitStatus, error) {
	// Get unit state.
	unitState, err := this.unitState(name)
	if err != nil {
		return UnitStatus{}, errgo.Mask(err)
	}

	// Get machine ip.
	ip, err := this.getMachineIp(name)
	if err != nil {
		return UnitStatus{}, errgo.Mask(err)
	}

	// Get unit.
	u, err := this.client.Unit(name)
	if err != nil {
		return UnitStatus{}, errgo.Mask(err)
	}

	return UnitStatus{
		Unit:   u.Name,
		State:  string(u.TargetState),
		Load:   unitState.LoadState,
		Active: unitState.ActiveState,
		Sub:    unitState.SubState,

		Description: u.Unit.Description(),

		Machine: ip,
	}, nil
}

func (this *ClientAPI) unitState(unitName string) (unit.UnitState, error) {
	unitStates, err := this.client.UnitStates()
	if err != nil {
		return unit.UnitState{}, errgo.Mask(err)
	}

	state := &unit.UnitState{}
	for _, unitState := range unitStates {
		if unitState.UnitName == unitName {
			state = unitState
		}
	}

	return *state, nil
}

func (this *ClientAPI) getMachineIp(unitName string) (string, error) {
	su, err := this.client.ScheduledUnit(unitName)
	if err != nil {
		return "", errgo.Mask(err)
	}

	machines, err := this.client.Machines()
	if err != nil {
		return "", errgo.Mask(err)
	}

	ip := ""
	for _, machine := range machines {
		if machine.ID == su.TargetMachineID {
			ip = machine.PublicIP
		}
	}

	return GetMachineIP(ip), nil
}
