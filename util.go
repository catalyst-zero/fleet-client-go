package client

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	execPkg "os/exec"
)

func isRunning(unitFileName, listUnitsOutPut string) (bool, error) {
	running := false
	scanner := bufio.NewScanner(strings.NewReader(listUnitsOutPut))

	// Scan each line of input.
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, unitFileName) {
			continue
		}

		for _, word := range strings.Split(line, "\t") {
			if word == "running" {
				running = true
			}
		}
	}

	// When finished scanning if any error other than io.EOF occured
	// it will be returned by scanner.Err().
	if err := scanner.Err(); err != nil {
		return false, scanner.Err()
	}

	return running, nil
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
