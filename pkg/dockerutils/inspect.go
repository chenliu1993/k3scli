package dockerutils

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"

	log "github.com/sirupsen/logrus"
)

// inspect container inspect container's infomation
func (c *ContainerCmd) Inspect() (string, error) {
	args := c.Args
	args = append(args, c.ID)
	cmd := exec.Command(c.Command, args...)
	info := ""
	infoBuf := bytes.NewBufferString(info)
	cmd.Stdout = infoBuf
	cmd.Stderr = os.Stderr
	log.Debug(fmt.Sprintf("begin inspect container %s", c.ID))
	err := cmd.Run()
	if err != nil {
		log.Debug(err)
		return "", err
	}
	return infoBuf.String(), nil
}