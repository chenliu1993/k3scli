# k3scli

## Introduction:
	k3scli is a command tool based on docker like kind(kubernetes-in-docker) but uses images with k3s in it.
	Currently it is just a tool under developing.

## Build
	```console
	go build main.go.
	```

## Usage
	To create a cluster, just use:
	```console
	go run main.go create <--config your-config.yaml> <cluster-name>
	```	
	There is an example config.yaml in the repo.
	If not specified, a default three-node cluster config will be used.

	To delete a cluster:
	```console
	go run main.go delete <cluster-name>
	```
	This should delete all resources on host and kill containers.	

	There are also some container-level functions like run/attach/kill a container, but it is just a wrap-use of docker with some options set to k3s's default values.

## TODO:
	1. support multiservers.
	2. worker roles show-up.
	3. pack neccessary images in tar format like kind.
## STUCK:	
	China's GFW.
