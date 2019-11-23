package dockerutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Kill wraps docker kill command
// and delete all related files on host
func (c *ContainerCmd) Kill(signal string) error {
	args := []string{
		"kill",
	}
	args = append(args,
		"--signal", signal,
		 c.ID,
		)
	cmd := exec.Command(c.Command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debug(fmt.Sprintf("begin run container %s", c.ID))
	err := cmd.Run()
	if err != nil {
		log.Debug(err)
		return err
	}
	// Delete all files under /tmp/k3s/<containerID>
	if err := os.RemoveAll(filepath.Join(K3sServerFile, c.ID)); err != nil {
		log.Debug(err)
		return err
	}
	if err := os.RemoveAll(filepath.Join(KubeCfgFolder, c.ID)); err != nil {
		log.Debug(err)
		return err
	}
	return nil
}