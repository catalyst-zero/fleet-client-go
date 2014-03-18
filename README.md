# fleet-client-go
Currently that client requires a `fleetctl` binary available to execute. Thus
there is no ssh connection used, the possible commands are `submit`, `start`,
`stop`, `destroy`, `list-units`, `list-machines`. The client supports a dump
version of `Status`, that uses `list-units` and just parses if a service is
running or not.

## install
```go
import fleetClientPkg "github.com/catalyst-zero/fleet-client-go"
```

## usage
```go
// Create new fleet client.
fleetClient := fleetClientPkg.NewClient()

// Submit unit file.
unitFilePath := "/tmp/unit-files/hello-world.service"
err := fleetClient.Submit(unitFilePath)

// Start a unit.
unitFileName := "hello-world.service"
err := fleetClient.Start(unitFileName)

// Stop a unit.
unitFileName := "hello-world.service"
err := fleetClient.Stop(unitFileName)

// Start a unit.
unitFileName := "hello-world.service"
err := fleetClient.Start(unitFileName)
