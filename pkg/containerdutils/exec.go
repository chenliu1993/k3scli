package containerdutils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
)

// Exec uses containerd exec to actually exec a process in a container
func (c *ContainerCmd) Exec(in io.Reader, out, stderr io.Writer) error {
	args := []string{
		"tasks",
		"exec",
	}
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
	err := cmd.Run()
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}