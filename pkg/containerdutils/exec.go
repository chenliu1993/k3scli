package containerdutils

import (
	"strconv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

var (
	// ExecID used for storing --exec-id
	ExecID = 1
)

// Exec uses containerd exec to actually exec a process in a container
func (c *ContainerCmd) Exec() error {
	args := []string{
		"tasks",
		"exec",
	}
	ExecID = ExecID + 1
	args = append(args, "--exec-id",
				strconv.Itoa(ExecID))
	args = append(args, c.ID)
	args = append(args, c.Args...)
	cmd := exec.Command(c.Command, args...)
	if c.Stdout != nil {
		cmd.Stdout = c.Stdout
	}
	if c.Stderr != nil {
		cmd.Stderr = c.Stderr
	}
	if c.Stdin != nil {
		cmd.Stdin = c.Stdin
	}
	log.Debug(fmt.Sprintf("begin exec process in container: %s", c.ID))
	fmt.Print(args)
	err := cmd.Run()
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}