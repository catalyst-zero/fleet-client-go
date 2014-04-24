package client

import (
	"bytes"
	"fmt"

	execPkg "os/exec"
)

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

func filterNonEmpty(values []string) []string {
	result := make([]string, 0)
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			result = append(result, v)
		}
	}
	return result
}
