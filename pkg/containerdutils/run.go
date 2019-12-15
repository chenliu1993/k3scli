package containerdutils

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// Run for running a container by containerd
func (c *ContainerCmd) Run() error {
	args := []string{
		"run",
		"--privileged",
		"--no-pivot",
		"--detach",
	}
	ctrFiles := filepath.Join(K3sServerFile, c.ID)
	if err := checkDir(ctrFiles); err != nil {
		return fmt.Errorf("kubeserver path failed")
	}
	ctrCfg := filepath.Join(KubeCfgFolder, c.ID)
	if err := checkDir(ctrCfg); err != nil {
		return fmt.Errorf("kubeconfig path failed")
	}
	args = append(args,
		"--env", "K3S_KUBECONFIG_OUTPUT="+filepath.Join(ctrCfg, "kubeconfig.yaml"),
		"--env", "K3S_KUBECONFIG_MODE=666",
		"--mount", "type=bind,src=/lib/modules,dst=/lib/modules,options=rbind:ro",
		"--mount", "type=bind,src=/var/lib/rancher/k3s,dst=ctrFiles,options=rbind:rw",
	)
	args = append(args, c.Args...)
	args = append(args, K3sBaseImage)
	if c.ID != "" {
		args = append(args,
			c.ID)
	}
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
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}