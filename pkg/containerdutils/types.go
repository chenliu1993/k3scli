package containerdutils

import (
	"fmt"
	"os"
	"io"
)

const (
	// KubeCfgFolder used as volume folder on host to store server config.
	KubeCfgFolder = "/tmp/k3s/configs"
	// K3sServerFile used as volume folder to store serve files like token.
	K3sServerFile = "/tmp/k3s/files"
	// K3sBaseImage for k3scli usage.
	K3sBaseImage = "docker.io/cliu2/k3sbase:0.11"
)

// ContainerCmd in containerdutils are using for wrapping containerd ctr.
type ContainerCmd struct {
	ID string
	Command string
	Args []string
	Stdout io.Writer
	Stderr io.Writer
	Stdin io.Reader
}

// SetStderr sets err
func (c *ContainerCmd) SetStderr(w io.Writer) {
	c.Stderr = w
}

// SetStdout sets out
func (c *ContainerCmd) SetStdout(w io.Writer) {
	c.Stdout = w
}
// SetStdin sets i
func (c *ContainerCmd) SetStdin(r io.Reader) {
	c.Stdin = r
}

func checkDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("remove existing k3s files")
		}
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("create k3s files failed")
		}
	}
	return nil
}