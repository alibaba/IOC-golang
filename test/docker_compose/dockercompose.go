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

package docker_compose

import (
	"log"
	"os/exec"
	"time"
)

const waitTime = time.Second * 3

func DockerComposeUp(path string, addtionalWaitTime time.Duration) error {
	out, err := exec.Command(
		"docker-compose",
		"-f", path,
		"up", "-d",
		"--remove-orphans").CombinedOutput()
	log.Printf("%s\n", string(out))
	time.Sleep(waitTime)
	time.Sleep(addtionalWaitTime)
	return err
}

func DockerComposeDown(path string) error {
	out, err := exec.Command(
		"docker-compose",
		"-f", path,
		"down", "-v").CombinedOutput()
	log.Printf("%s\n", string(out))
	time.Sleep(waitTime)
	return err
}
