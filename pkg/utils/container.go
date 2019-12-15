package utils

import (
	"fmt"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
	containerd "github.com/chenliu1993/k3scli/pkg/containerdutils"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	log "github.com/sirupsen/logrus"
)

const (
	// NODE_VERSION = "0.10"
	// NODE_VERSION = "withimages"
	// NODE_IMAGE   = "cliu2/k3snode:" + NODE_VERSION

	// BaseVersion defines the version of k3sbase image.
	BaseVersion = "0.11"
	// BaseImage defines the image used by k3scli.
	BaseImage = "cliu2/k3sbase:" + BaseVersion
)

// RunContainer used for wrap exec run
func RunContainer(containerID string, label string, detach bool, image string, ports []string, cluster string, mode string) (err error) {
	log.Debug("generating run cmd")
	if mode == "docker" {
		ctrCmd := docker.ContainerCmd{
			ID:      containerID,
			Command: "docker",
		}
		ctrCmd.Args = []string{}
		for _, port := range ports {
			ctrCmd.Args = append(ctrCmd.Args, "-p", port)
		}
		if cluster != "" {
			ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("%s=%s", clusterconfig.ClusterLabelKey, cluster))
		}
		if label != "" {
			ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("%s=%s", clusterconfig.ClusterRole, cluster+"-"+label))
		}
		// for deploy pod to a specific node
		ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("type=%s", containerID))
		ctrCmd.Detach = detach
		ctrCmd.Image = image
		err = ctrCmd.Run()
	} else if mode == "containerd" {
		ctrCmd := containerd.ContainerCmd{
			ID: containerID,
			Command: "ctr",
		}
		ctrCmd.Args = []string{}
		if cluster != "" {
			ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("%s=%s", clusterconfig.ClusterLabelKey, cluster))
		}
		if label != "" {
			ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("%s=%s", clusterconfig.ClusterRole, cluster+"-"+label))
		}
		// for deploy pod to a specific node
		ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("type=%s", containerID))
		err = ctrCmd.Run()
	}
	return err
}

// Join used as join interface for a agent to join the server node.
func Join(containerID, server, token string, detach bool) error {
	log.Debug("generating exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	// Has to be true, because k3scli now it is not a input tty
	ctrCmd.Detach = detach
	// k3s agent --server https://myserver:6443 --token ${NODE_TOKEN}
	ctrCmd.Args = []string{
		"k3s", "agent",
		"--server", server,
		"--token", token,
	}
	fmt.Print(ctrCmd.Args)
	return ctrCmd.Exec(nil, nil, nil)
}

// AttachContainer attatches io to a container.
// warpper for docker exec.
func AttachContainer(containerID string) error {
	log.Debug("generating docker exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	// Has to be false, because attach means interact with container
	ctrCmd.Detach = false
	// just a sh command
	ctrCmd.Args = []string{
		"/bin/sh",
	}
	return ctrCmd.Exec(nil, nil, nil)
}

// KillContainer kills a container.
func KillContainer(containerID, signal string) error {
	log.Debug("generating docker exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	return ctrCmd.Kill(signal)
}

// InspectContainerIP returns the ip of a container.
func InspectContainerIP(containerID, mode string) (ip string, err error) {
	if mode == "docker" {
		ctrCmd := docker.ContainerCmd{
			ID:      containerID,
			Command: "docker",
			Args: []string{"inspect", "--format",
				"'{{.NetworkSettings.IPAddress}}'"},
		}
		ip, err = ctrCmd.Inspect()
	} else if mode == "containerd" {
		ctrCmd := containerd.ContainerCmd{
			ID: containerID,
			Command: "ctr",
			Args: []string{},
		}
		ip, err = GetCtrIP()
	}
	return ip, err
}

// ExecInContainer executes a command in the target container.
func ExecInContainer(containerID, cmd string, detach bool) (err error) {
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	ctrCmd.Detach = detach
	ctrCmd.Args = []string{
		"sh", "-c",
		cmd,
	}
	// fmt.Print(ctrCmd.Args)
	return ctrCmd.Exec(nil, nil, nil)
}

// GetCtrIP is the containerd-version of get server ip
func GetCtrIP() (ip string, err error) {
	
}