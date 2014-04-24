package client

import (
	"bytes"
	"fmt"
	"strings"

	execPkg "os/exec"
)

// GetMachineIP parses the unitMachine in format "uuid/ip" and returns only the IP part.
// Can be us with the {UnitStatus.Machine} field.
// Returns an empty string, if no ip was found.
func GetMachineIP(unitMachine string) string {
	fields := strings.Split(unitMachine, "/")
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}

func exec(cmd *execPkg.Cmd) (string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	if err := stderr.String(); err != "" {
		return "", fmt.Errorf(err)
	}

	return stdout.String(), nil
}

func filterEmpty(values []string) []string {
	result := make([]string, 0)
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			result = append(result, v)
		}
	}
	return result
}
