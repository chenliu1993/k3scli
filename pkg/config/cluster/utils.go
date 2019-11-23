package config

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)

// Load loads a cluster config from a yaml file
func Load(config string) (Cluster, error) {
	cluster := Cluster{}
	log.Debug("loading cluster config")
	yamlData, err := ioutil.ReadFile(config)
	if err != nil {
		log.Debug(err)
		return cluster, err
	}
	err = yaml.Unmarshal(yamlData, &cluster)
	if err != nil {
		log.Debug(err)
		return cluster, err
	}
	if cluster.Name == "" {
		cluster.Name = DefaultClusterName
	}
	return cluster, err
}