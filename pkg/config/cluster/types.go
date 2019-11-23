package config

import (

)

const (
	ClusterLabelKey = "kind.k3s.cluster"
	ClusterRole = "kind.k3s.cluster.role"
)

const (
	// DefaultClusterName 
	// when yaml does not points out a cluster name,
	// this will be used
	DefaultClusterName = "k3scluster"
)

// Cluster represents a k3s container cluster
type Cluster struct {
	Name string `yaml:"cluster_name"`
	Nodes []Node `yaml:"nodes"`
}

// Node represents a cluster k3s node
type Node struct {
	Name string
	Label string `yaml:"role"`
}




