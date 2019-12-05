package utils

import (
	"fmt"

	// "bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultArchivesPath   = "/k3s"
	DefaultContainerdSock = "/run/k3s/containerd/containerd.sock"
)

// This file contains some node related functions like get server's ip and token

// GetServerToken get server token content
func GetServerToken(containerID string) (string, error) {
	// time.Sleep(15 * time.Second)
	// First copy token out from container
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
		Args:    []string{"cp"},
	}
	ctrCmd.Detach = true
	ctrCmd.Args = append(ctrCmd.Args,
		containerID+":"+docker.K3sServerFileInContainer,
		filepath.Join(docker.K3sServerFile, containerID),
	)
	cmd := exec.Command(
		ctrCmd.Command, ctrCmd.Args...,
	)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	// Second read from token file
	bytes, err := ioutil.ReadFile(filepath.Join(docker.K3sServerFile, containerID, "token"))
	if err != nil {
		log.Debug(err)
		return "", err
	}
	tokenStr := strings.Replace(string(bytes), "\n", "", -1)
	return strings.TrimSpace(string(tokenStr)), nil
}

// GetServerIP get server internal IP through docker inspect
func GetServerIP(containerID string) (string, error) {
	log.Debug("get server ip from docker inspect")
	ip, err := InspectContainerIP(containerID)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	// remove the unneccessary '
	ip = ip[1 : len(ip)-2]
	server := "https://" + ip + ":6443"
	return server, nil
}

// GetClusterNames and returns it
func GetClusterNames(clusterName string) (lines []string, err error) {
	// For now, only supports one server, so server name will be based on th cluster name
	cmd := exec.Command(
		"docker",
		"ps",
		"-q",         // quiet output for parsing
		"-a",         // show stopped nodes
		"--no-trunc", // don't truncate
		// filter for nodes with the cluster label
		"--filter", fmt.Sprintf("label=%s=%s", clusterconfig.ClusterLabelKey, clusterName),
		// format to include the cluster name
		"--format", `{{.Names}}`,
	)
	lines, err = docker.ExecOutput(*cmd, false)
	if err != nil {
		return nil, err
	}
	// currentlt only supports one server
	// if len(lines) != 1 {
	// 	return nil, fmt.Errorf("k3scli don't support multiserver now...")
	// }
	return lines, nil
}

// Generate container a unique container name
func GenCtrName() string {
	return uuid.New().String()
}

// LoadImages used for load images in tar
func LoadImages(containerID string) error {
	//for mac just load pause.tar is ok
	findCmd := "find " + DefaultArchivesPath + " -name pause.tar"
	loadCmd := "xargs -n1 k3s ctr -a " + DefaultContainerdSock + " images import"
	cmd := findCmd + " | " + loadCmd
	return ExecInContainer(containerID, cmd, false)
}

// StartK3S	starts k3s daemon service
func StartK3S(containerID string) error {
	log.Debug("starting k3s server")
	startCmd := "nohup k3s server"
	err := ExecInContainer(containerID, startCmd, true)
	if err != nil {
		return err
	}
	return nil
}

// GenratePortMapping
func GenratePortMapping(ports []clusterconfig.Port) []string {
	var portmappings []string
	for _, port := range ports {
		portmapping := port.Hostport + ":" + port.Port
		portmappings = append(portmappings, portmapping)
	}
	return portmappings
}

// CopyFromHostToCtr copies a file into the container
// follow Kind
func CopyFromHostToCtr(containerID, file string) (err error) {
	log.Debug("get file's content into buffer...")
	// var buff bytes.Buffer
	current, err := os.Getwd()
	if err != nil {
		return err
	}
	// first copy host file's content into buffer
	// cmd := exec.Command(
	// 	"cat", filepath.Join(current, file),
	// )
	// cmd.Stdout = &buff
	// if err = cmd.Run(); err != nil {
	// 	return err
	// }
	// then copy from buffer to container with the same name
	ctrCmd := docker.ContainerCmd{
		ID:      containerID,
		Command: "docker",
		Args:    []string{"cp"},
	}
	ctrCmd.Detach = true
	ctrCmd.Args = append(ctrCmd.Args,
		filepath.Join(current, file),
		containerID+":/"+file,
	)
	cmd := exec.Command(
		ctrCmd.Command, ctrCmd.Args...,
	)
	return cmd.Run()
	// ctrCmd.Args = []string{
	// 	"cp", "/dev/stdin",
	// 	file,
	// }
	// ctrCmd.Detach = true
	// if err = ctrCmd.Exec(&buff,nil,nil); err != nil {
	// 	return err
	// }
	// return nil
}
