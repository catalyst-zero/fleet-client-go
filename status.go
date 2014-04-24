package client

// Parse for `fleet status` output
import (
	"bufio"
	"fmt"
	execPkg "os/exec"
	"strings"
)

const (
	LOAD_UNKNOWN = "-"
	LOAD_LOADED  = "loaded"

	ACTIVE_UNKNOWN = "-"
	ACTIVE_ACTIVE  = "-"
	ACTIVE_FAILED  = "failed"

	SUB_UNKNOWN = "-"
	SUB_RUNNING = "running"
	SUB_FAILED  = "failed"
)

type UnitStatus struct {
	// Unit Name with file extension
	Unit string

	// "-", "loaded"
	Load string

	// "-", "active", "failed"
	Active string

	// The state of the unit, e.g. "-", "running" or "failed". See the SUB_* constants.
	Sub         string
	Description string

	// The machine that is used to execute the unit.
	// Is "-", when no machine is assigned.
	// Otherwise is in the format of "uuid/ip", where uuid is shortened version of the host uuid
	// and IP is the IP assigned to that machine.
	Machine string
}

// MachineIP returns the IP of the used machine, if available, an empty string otherwise.
// Alias for `return GetMachineIp(status.Machine)`.
func (status UnitStatus) MachineIP() string {
	return GetMachineIp(status.Machine)
}

// StatusAll executes "fleetctl status" and parses the output table. Thus, certain fields can be mangled or
// shortened, e.g. the machine column.
func (this *Client) StatusAll() ([]UnitStatus, error) {
	cmd := execPkg.Command(FLEETCTL, ENDPOINT_OPTION, this.etcdPeer, "list-units")
	stdout, err := exec(cmd)
	if err != nil {
		return []UnitStatus{}, err
	}

	return parseFleetStatusOutput(stdout)
}

func parseFleetStatusOutput(output string) ([]UnitStatus, error) {
	result := make([]UnitStatus, 0)

	scanner := bufio.NewScanner(strings.NewReader(output))
	// Scan each line of input.
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		if lineCount == 1 {
			continue
		}

		words := filterEmpty(strings.Split(line, "\t"))
		unitStatus := UnitStatus{
			Unit:        words[0],
			Load:        words[1],
			Active:      words[2],
			Sub:         words[3],
			Description: words[4],
			Machine:     words[5],
		}
		result = append(result, unitStatus)
	}

	// When finished scanning if any error other than io.EOF occured
	// it will be returned by scanner.Err().
	if err := scanner.Err(); err != nil {
		return result, scanner.Err()
	}
	return result, nil
}

// StatusUnit returns the UnitStatus for the given unitfile.
// If the unit is not found or could not be retrieved, an error will be returned.
//
// Internally it executes `fleetctl status` and looks for the matching
// unit file. Its {UnitStatus} will be returned.
func (this *Client) StatusUnit(unitFileName string) (UnitStatus, error) {
	status, err := this.StatusAll()
	if err != nil {
		return UnitStatus{}, err
	}

	for _, s := range status {
		if s.Unit == unitFileName {
			return s, nil
		}
	}
	return UnitStatus{}, fmt.Errorf("Unknown unitfilename: %s", unitFileName)
}

type Status struct {
	Running     bool
	ContainerIP string
}

func (this *Client) Status(unitFileName string) (Status, error) {
	allStatus, err := this.StatusAll()
	if err != nil {
		return Status{}, err
	}

	for _, status := range allStatus {
		if status.Unit == unitFileName {
			result := Status{
				Running:     status.Sub == SUB_RUNNING,
				ContainerIP: "127.0.0.1",
			}
			return result, nil
		}
	}
	// Return running=false, because we didn't find it
	return Status{}, nil
}
