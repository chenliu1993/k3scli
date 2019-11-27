package utils

import (
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
)


// CreateCluster creates a cluster given the default name
func CreateCluster(clusterName string, cluster clusterconfig.Cluster, ports []string) (err error) {
	log.Debug("Creating cluster...")
	var name string
	// First running server container
	serverName := GenCtrName()
	if ports == nil {
		ports = []string{}
	}
	// deal with port mapping
	originPorts, err := ConvertFromStrToInt(ports)
	if err != nil {
		return err
	}
	serverPorts := append(ports, "6443")
	err = RunContainer(serverName, "server", true, BASE_IMAGE, serverPorts, clusterName)
	if err != nil {
		return err
	}
	err = StartK3S(serverName)
	if err != nil {
		return err
	}
	time.Sleep(2*time.Second)
	err = LoadImages(serverName) 
	if err != nil {
		return err
	}
	cluster.Nodes[0].Name = serverName
	server, err := GetServerIP(serverName)
	if err != nil {
		return err
	}
	serverToken, err := GetServerToken(serverName)
	if err != nil {
		return err
	}
	var newPorts = originPorts
	// Second join worker nodes one-by-one
	// Join(containerID, server, token, detach)
	// Server node must on the first place of config file
	for _, node := range cluster.Nodes[1:] {
		// First run container then join container
		name = GenCtrName()
		newPorts = AddPort(newPorts)
		fmt.Print(newPorts)
		err := RunContainer(name, "worker", true, BASE_IMAGE, ConvertFromIntToStr(newPorts), clusterName)	
		if err != nil {
			return err
		}
		node.Name = name
		if err := Join(name, server, serverToken, true); err != nil {
			return err
		}
		time.Sleep(3*time.Second)
		err = LoadImages(name) 
		if err != nil {
			return err
		}
	}
	
	return nil
}

// DeleteCluster first gets all cluster's container names, then kill them one-by-one
func DeleteCluster(clusterName string) error {
	log.Debugf("get all containers under %s", clusterName)
	names, err := GetClusterNames(clusterName)
	if err != nil {
		return err
	}
	for _, name := range names {	
		log.Debugf("killing container: %s", name)
		err := KillContainer(name, "sigterm")
		if err != nil {
			return err
		}
	}
	return nil
}