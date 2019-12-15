package utils

import (
	"fmt"
	"time"
	// "bytes"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// DefaultArchivesPath is the path in container where image tars are stored.
	DefaultArchivesPath = "/k3s"
	// DefaultContainerdSock is the default containerd socket file in container.
	DefaultContainerdSock = "/run/k3s/containerd/containerd.sock"
)

// This file contains some node related functions like get server's ip and token

// GetServerToken get server token content
func GetServerToken(containerID string) (string, error) {
	log.Debug("read token out from k3s server files")
	time.Sleep(10 * time.Second)
	// fileInfo, err := os.Stat(filepath.Join(docker.K3sServerFile, containerID, "server"))
	// if err != nil {
	// 	fmt.Print(err)
	// 	return "", err
	// }
	// fmt.Print(fileInfo.Name())
	// token place
	token := filepath.Join(docker.K3sServerFile, containerID, "server", "token")
	bytes, err := ioutil.ReadFile(token)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	tokenStr := strings.Replace(string(bytes), "\n", "", -1)
	return strings.TrimSpace(string(tokenStr)), nil
}

// GetServerIP get server internal IP through docker inspect
func GetServerIP(containerID, mode string) (server string, err error) {
	log.Debug("get server ip")
	ip, err := InspectContainerIP(containerID, mode)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	// remove the unneccessary '
	ip = ip[1 : len(ip)-2]
	server = "https://" + ip + ":6443"	
	
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

// GenCtrName generate container a unique container name
func GenCtrName() string {
	return uuid.New().String()
}

// LoadImages use ctr to load images that is in the form of tar
func LoadImages(containerID string, role string) error {
	log.Debug("loading images")
	var findCmd string
	// list image tars
	if role == "server" {
		findCmd = "find " + DefaultArchivesPath + " -name *.tar"
	} else if role == "worker" {
		findCmd = "find " + DefaultArchivesPath + " -name *-lb.tar"
	}
	loadCmd := "xargs -n1 k3s ctr -a " + DefaultContainerdSock + " images import"
	Cmd := findCmd + " | " + loadCmd
	err := ExecInContainer(containerID, Cmd, false)
	if err != nil {
		return err
	}
	if role == "worker" {
		findCmd = "find " + DefaultArchivesPath + " -name pause.tar"
		Cmd := findCmd + " | " + loadCmd
		err := ExecInContainer(containerID, Cmd, false)
		if err != nil {
			return err
		}
		findCmd = "find " + DefaultArchivesPath + " -name *traefik*.tar"
		Cmd = findCmd + " | " + loadCmd
		err = ExecInContainer(containerID, Cmd, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// StartK3S starts k3s daemon service
func StartK3S(containerID string) error {
	log.Debug("starting k3s server")
	startCmd := "nohup k3s server"
	err := ExecInContainer(containerID, startCmd, true)
	if err != nil {
		return err
	}
	return nil
}

// GenratePortMapping takes inpout from config
// and generates []string
// with each component is a port mapping pair like:
// -p 9000:9000
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
}
