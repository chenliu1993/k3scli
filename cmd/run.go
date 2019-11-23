package cmd

import (
	"strings"
	"context"

	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// RunCommand wraps docker run for k3scli
var RunCommand = cli.Command{
	Name:  "run",
	Usage: "run a k3sbase/k3snode container",
	ArgsUsage: `<container-id> is your name for the instance of the container that you
	are starting. The name you provide for the container instance must be unique
	on your host.`,
	Description: `The run command allows you to start a new k3sbase/k3snode container`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "label, l",
			Value: "",
			Usage: `label used for docker run --label used for distinguishing from server to worker (primary)`,
		},
		&cli.BoolFlag{
			Name:  "detach, d",
			Usage: `run in detach mode or not`,
		},
		&cli.StringFlag{
			Name:  "image",
			Usage: `image used`,
		},
		&cli.StringSliceFlag{
			Name:  "port, p",
			Usage: `port mapping between container and host`,
		},
		&cli.StringFlag{
			Name: "cluster",
			Usage: `cluster this contaienr belongs to`,
		},
	},
	Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return run(ctx, context.Args().First(),
			context.String("label"),
			context.Bool("detach"),
			context.String("image"),
			context.StringSlice("port"),
			context.String("cluster"),
		)
	},
}

func run(ctx context.Context, containerID, label string, detach bool, image string, ports []string, cluster string) error {
	log.Debug("begin running container")
	if label == "" {
		log.Debug("role of container is not set, default to server")
		label = "server"
	}
	if image == "" {
		log.Debug("k3s image not set, default to node")
		image = utils.NODE_IMAGE
	}
	if label == "server" && strings.Index(image, "base") != -1 {
		log.Fatal("base image cannot serve as server")
	}
	if label == "worker" && strings.Index(image, "node") != -1 {
		log.Fatal("node image cannot serve as worker")
	}
	if containerID == "" {
		log.Debug("no container ID, will generate a random one by docker")
	}
	if ports == nil {
		if label == "server" {
			ports = append(ports, "6443")
		} else {
			log.Debug("no ports mapping is defined")
			ports = []string{}
		}
	}
	if cluster == "" {
		log.Debug("no cluster specified")
	}
	return utils.RunContainer(containerID, label, detach, image, ports, cluster)
}
