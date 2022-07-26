package iocli_command

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

const whichCommand = "which"
const iocliCommand = "iocli"

func getCommandPath(commandName string) (string, error) {
	out, err := exec.Command(whichCommand, commandName).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

func Run(command []string, timeout time.Duration) (string, error) {
	iocliBinPath, err := getCommandPath(iocliCommand)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	cmd := exec.Command(
		iocliBinPath,
		command...)
	cmd.Stdout = &b
	cmd.Stderr = &b
	closeCh := make(chan struct{})
	go func() {
		_ = cmd.Run()
		close(closeCh)
	}()

	after := time.After(timeout)
	select {
	case <-closeCh:
	case <-after:
	}
	return b.String(), nil
}
