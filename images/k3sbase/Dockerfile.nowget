
FROM alpine:3.10

MAINTAINER LIUCHEN
COPY files/ /
COPY temp/ /

RUN apk update; \
	apk add --no-cache --virtual .build-deps \
        	iptables iproute2 ethtool socat util-linux ebtables udev kmod \
     		bash ca-certificates curl rsync \ 
		containerd docker openrc wget tar; \
	mkdir -p /run/openrc; \
	touch /run/openrc/softlevel; \
	export PATH=$PATH:/usr/local/bin; \
	tar -C /usr/local/bin -xzvf crictl-v1.16.1-linux-amd64.tar.gz; \
	rm -rf crictl-v1.16.1-linux-amd64.tar.gz; \
	mkdir -p /opt/cni/bin; \
	tar -C /opt/cni/bin -xzvf cni-plugins-linux-amd64-v0.8.3.tgz; \
	rm -rf cni-plugins-linux-amd64-v0.8.3.tgz; \
	mkdir -p /etc/kubernetes/manifests;

ENV container docker

ENTRYPOINT [ "/bin/sh" ] 
