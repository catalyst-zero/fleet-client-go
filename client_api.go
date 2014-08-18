package client

import (
	"fmt"
	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/schema"
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

	c, _ := client.NewHTTPClient(&hc, "http://localhost/")

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
	uf, err := getUnitFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed reading unit file %s: %v", name, err)
	}

	options := schema.MapUnitFileToSchemaUnitOptions(uf)

	unit := schema.Unit{
		Name:    name,
		Options: options,
	}

	if err := this.client.CreateUnit(&unit); err != nil {
		return fmt.Errorf("failed creating unit %s: %v", unit.Name, err)
	}

	return nil
}

func (this *ClientAPI) Unit(name string) (*schema.Unit, error) {
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
	this.client.SetUnitTargetState(u.Name, STATE_LAUNCHED)

	return nil
}

func (this *ClientAPI) Stop(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetUnitTargetState(u.Name, STATE_LOADED)

	return nil
}

func (this *ClientAPI) Load(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.SetUnitTargetState(u.Name, STATE_LOADED)

	return nil
}

func (this *ClientAPI) Destroy(name string) error {
	u, err := this.Unit(name)

	if err != nil {
		return errgo.Mask(err)
	}
	this.client.DestroyUnit(u.Name)

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

	// Get unit.
	u, err := this.client.Unit(name)
	if err != nil {
		return UnitStatus{}, errgo.Mask(err)
	}

	// Get machine ip.
	ip, err := this.getMachineIp(u)
	if err != nil {
		return UnitStatus{}, errgo.Mask(err)
	}

	return UnitStatus{
		Unit:        u.Name,
		Description: description(u.Options),

		State:  string(u.DesiredState),
		Load:   unitState.SystemdLoadState,
		Active: unitState.SystemdActiveState,
		Sub:    unitState.SystemdSubState,

		Machine: ip,
	}, nil
}

func (this *ClientAPI) unitState(unitName string) (schema.UnitState, error) {
	unitStates, err := this.client.UnitStates()
	if err != nil {
		return schema.UnitState{}, errgo.Mask(err)
	}

	state := &schema.UnitState{}
	for _, unitState := range unitStates {
		if unitState.Name == unitName {
			state = unitState
		}
	}

	return *state, nil
}

func (this *ClientAPI) getMachineIp(unit *schema.Unit) (string, error) {
	machines, err := this.client.Machines()
	if err != nil {
		return "", errgo.Mask(err)
	}

	ip := ""
	for _, machine := range machines {
		if machine.ID == unit.MachineID {
			ip = machine.PublicIP
		}
	}

	return GetMachineIP(ip), nil
}

func description(options []*schema.UnitOption) string {
	for _, option := range options {
		if option.Section == "Unit" && option.Name == "Description" {
			return option.Value
		}
	}

	return ""
}
