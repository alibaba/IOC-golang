/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
	return b.String(), nil
}
