package client

import (
	"bufio"
	"strings"
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

		for _, word := range strings.Split(line, " ") {
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
