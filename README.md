# fleet-client-go

Currently that client requires a `fleetctl` binary available to execute. Thus
there is no ssh connection used, the possible commands are `submit`, `start`,
`stop`, `destroy`, `list-units`, `list-machines`. The client supports a dump
version of `Status`, that uses `list-units` and just parses whether a service
is running or not.

A new way to use fleets http api has been implemented.

## install
```go
import fleetClientPkg "github.com/catalyst-zero/fleet-client-go"
```

## usage
```go
// Create new fleet client based on a given binary.
fleetClient := fleetClientPkg.NewClientCLI()

// Create new fleet client based on the http api.
fleetClient := fleetClientPkg.NewClientAPI()

// Interface methods.
type FleetClient interface {
  // A Unit is a submitted job known by fleet, but not started yet. Submitting
  // a job creates a unit. Unit() returns such an object. Further a Unit has
  // different properties than a ScheduledUnit.
	Unit(name string) (*job.Unit, error)

  // A ScheduledUnit is a submitted job known by fleet in a specific state.
  // ScheduledUnit() does not fetch a ScheduledUnit if a Unit is not started
  // yet, but only submitted. Further a ScheduledUnit has different properties
  // than a Unit.
	ScheduledUnit(name string) (*job.ScheduledUnit, error)

	Submit(name, filePath string) error
	Start(name string) error
	Stop(name string) error
	Load(name string) error
	Destroy(name string) error
	Status(name string) (*Status, error) // Deprecated, use StatusUnit()
	StatusUnit(name string) (UnitStatus, error)
}
```
