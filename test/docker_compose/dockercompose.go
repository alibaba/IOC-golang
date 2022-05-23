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
