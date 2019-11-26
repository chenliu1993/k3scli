
FROM alpine:3.10

MAINTAINER LIUCHEN

COPY files/ /
RUN mkdir -p /k3s
COPY images/ /k3s

RUN apk update; \
	apk add --no-cache --virtual .build-deps \
        	iptables iproute2 ethtool socat util-linux ebtables udev kmod \
     		bash ca-certificates curl rsync \ 
		containerd openrc wget tar; \
	mkdir -p /run/openrc; \
	touch /run/openrc/softlevel; \
	export PATH=$PATH:/usr/local/bin; \
	wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.16.1/crictl-v1.16.1-linux-amd64.tar.gz; \
	tar -C /usr/local/bin -xzvf crictl-v1.16.1-linux-amd64.tar.gz; \
	rm -rf crictl-v1.16.1-linux-amd64.tar.gz; \
	wget https://github.com/containernetworking/plugins/releases/download/v0.8.3/cni-plugins-linux-amd64-v0.8.3.tgz; \
	mkdir -p /opt/cni/bin; \
	tar -C /opt/cni/bin -xzvf cni-plugins-linux-amd64-v0.8.3.tgz; \
	rm -rf cni-plugins-linux-amd64-v0.8.3.tgz; \
	mkdir -p /etc/kubernetes/manifests; 

VOLUME ["/var/lib/rancher/k3s"]
	
ENV container docker

ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ] 
