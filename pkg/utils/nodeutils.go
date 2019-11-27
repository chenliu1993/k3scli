package utils

import (
	"fmt"
	"time"
	"strconv"
	"github.com/google/uuid"
	// "os"
	"os/exec"
	"strings"
	"path/filepath"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	docker "github.com/chenliu1993/k3scli/pkg/dockerutils"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
)


const (
	DefaultArchivesPath = "/k3s"
	DefaultContainerdSock = "/run/k3s/containerd/containerd.sock"
)
// This file contains some node related functions like get server's ip and token

// func Init() {
// 	file, err := os.Open("/dev/urandom")
//     if err != nil {
//             panic(fmt.Sprintf("Failed to open urandom: %v", err))
//     }
//     uuid.SetRand(file)
// }
// GetServerToken get server token content
func GetServerToken(containerID string) (string, error) {
	log.Debug("read token out from k3s server files")
	time.Sleep(10*time.Second)
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
func GetServerIP(containerID string) (string, error) {
	log.Debug("get server ip from docker inspect")
	ip, err := InspectContainerIP(containerID)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	// remove the unneccessary '
	ip = ip[1:len(ip)-2]
	server := "https://"+ip+":6443"
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

// LoadImages use ctr to load images that is in the form of tar
func LoadImages(containerID string) error {
	log.Debug("loading images")
	// list image tars
	findCmd := "find "+DefaultArchivesPath+" -name *.tar"
	loadCmd := "xargs -n1 k3s ctr -a "+DefaultContainerdSock+" images import"
	Cmd := findCmd+" | "+loadCmd
	err := ExecInContainer(containerID, Cmd, false)
	if err != nil {
		return err
	}
	return nil
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

// ConvertPorts convert string to int ports

func ConvertFromStrToInt(strs []string) ([]int, error) {
	var intStrs []int
	for _, str := range strs {
		intStr, err :=  strconv.Atoi(str)
		if err != nil {
			return nil ,err
		}
		intStrs = append(intStrs, intStr)
	}
	return intStrs, nil
}

func ConvertFromIntToStr(intStrs []int) ([]string) {
	var strs []string
	for _, intStr := range intStrs {
		strs = append(strs, strconv.Itoa(intStr))
	}
	return strs
}

func AddPort(ports []int) []int {
	var newPorts []int
	for _, port := range ports {
		newPorts = append(newPorts, port+1)
	}
	return newPorts
}