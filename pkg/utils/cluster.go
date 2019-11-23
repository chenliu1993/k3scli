package utils

import (
	log "github.com/sirupsen/logrus"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
)


// CreateCluster creates a cluster given the default name
func CreateCluster(clusterName string, cluster clusterconfig.Cluster) (err error) {
	log.Debug("Creating cluster...")
	var name string
	// First running server container
	serverName := GenCtrName()
	err = RunContainer(serverName, "server", true, NODE_IMAGE, []string{"6443"}, clusterName)
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
	// Second join worker nodes one-by-one
	// Join(containerID, server, token, detach)
	// Server node must on the first place of config file

	for _, node := range cluster.Nodes[1:] {
		// First run container then join container
		name = GenCtrName()
		err := RunContainer(name, "worker", true, BASE_IMAGE, []string{}, clusterName)	
		if err != nil {
			return err
		}
		node.Name = name
		if err := Join(name, server, serverToken, true); err != nil {
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