package utils

import (
	"fmt"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
	log "github.com/sirupsen/logrus"
)
const (
	// NODE_VERSION = "0.10"
	NODE_VERSION = "allsetup"
	NODE_IMAGE = "cliu2/k3snode:"+NODE_VERSION

	BASE_VERSION = "0.10"
	BASE_IMAGE = "cliu2/k3sbase:"+BASE_VERSION
)

// RunServerContainer used for wrap exec run
func RunContainer(containerID string, label string, detach bool, image string, ports []string, cluster string) error {
	log.Debug("generating docker run cmd")
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
	}
	ctrCmd.Args = []string{}
	for _, port := range ports {
		portStr := port+":"+port
		ctrCmd.Args = append(ctrCmd.Args, "-p", portStr)
	}
	if cluster != "" {
		ctrCmd.Args =  append(ctrCmd.Args, "--label", fmt.Sprintf("%s=%s", clusterconfig.ClusterLabelKey, cluster))
	}
	if label != "" {
		ctrCmd.Args = append(ctrCmd.Args, "--label", fmt.Sprintf("%s=%s", clusterconfig.ClusterRole, cluster+"-"+label))
	}
	ctrCmd.Detach = detach
	ctrCmd.Image = image
	return ctrCmd.Run()
}

func Join(containerID, server, token string, detach bool) error {
	log.Debug("generating docker exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID: containerID,
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
	return ctrCmd.Exec()
}


func AttachContainer(containerID string) error {
	log.Debug("generating docker exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID: containerID,
		Command: "docker",
	}
	// Has to be false, because attach means interact with container
	ctrCmd.Detach = false
	// just a sh command
	ctrCmd.Args = []string{
		"/bin/sh",
	}
	return ctrCmd.Exec()
}


func KillContainer(containerID, signal string) error {
	log.Debug("generating docker exec cmd")
	ctrCmd := docker.ContainerCmd{
		ID: containerID,
		Command: "docker",
	}
	return ctrCmd.Kill(signal)
}


func InspectContainerIP(containerID string) (string, error) {
	ctrCmd := docker.ContainerCmd{
		ID :containerID,
		Command: "docker",
		Args: []string{"inspect", "--format",
			"'{{.NetworkSettings.IPAddress}}'"},
	}
	return ctrCmd.Inspect()
}
