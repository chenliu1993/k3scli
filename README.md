# k3scli

## Introduction:
	k3scli is a command tool based on docker like kind(kubernetes-in-docker) but uses images with k3s in it.
	Currently it is just a tool under developing.
	Faster than Kind, lighter than Kind.

## Build
	```console
	go build main.go.
	```
## Dependency
	Before k3snode/k3sbase:0.11, k3scli will use two different images to start up a container. k3snode will serve as server while k3sbase will server as worker. But k3s in k3snode depends heavily on k8s.gcr.io images(pause:3.1). Developers in China cannot pull this image which is really annoying. So k3sbase:0.11 is decided, with k3s needed images present in it. Each time start up a server, ctr in k3sbase:0.11 will load the image tars therefore avoiding GFW issue. 
	Besides now server and worker can be both started up based on k3sbase:0.11. Although volumes of image is enlarged, but I think it is worth it.

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
	2. worker roles show-up (k3s without container does not show up worker roles too).
	3. pack neccessary images in tar format like kind (done).
	4. remove docker support to use containerd?
	5. modify codes architecture and quality before adding new functions on master.
## STUCK:	
	China's GFW.
## Branches:
	1. master is the main branch.
	2. topic/chenliu1993/macos is the MacOS version.
	3. topic/chenliu1993/containerd is the container as runtime version (releaizing...) 


