package cmd

import (
	"context"

	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)
// Join command combines run and join.
// Used to join a k3s worker node to a server

var (
        // May needs perfecting
        defaultPorts = []string{}
)

var JoinCommand = cli.Command{
	Name:  "join",
        Usage: "join a k3sbase container to a existing a server",
        ArgsUsage: `join <--detach> --server <SERVER-IP> --token <TOKEN> to <worker-container-id> <server-container-id`,
        Description: `The join command run a k3sbase container and join it to an existing k3snode server container`,
	Flags: []cli.Flag{
                &cli.StringFlag{
                        Name:  "server-id",
                        Value: "",
                        Usage: `server container id`,
                },
                // &cli.StringFlag{
                //         Name:  "token, t",
                //         Usage: `server token resides in /var/lib/rancher/k3s/server/node-token on server container`,
                // },
                &cli.BoolFlag{
                        Name:  "detach, d",
                        Usage: `detach mode`,
                },
        },
        Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
                return join(ctx, context.Args().First(),
                        context.String("server-id"),
                        context.Bool("detach"),
		)
        },
}

func join(ctx context.Context, containerID, serverID string, detach bool) error {
        log.Debug("Begin join server node, first checking args")
        // First run a worker container
        log.Debug("run worker container")
        // Detach has to be true, other wise the join action cannot execute.
        err := run(ctx, containerID, "worker", true, utils.BASE_IMAGE, defaultPorts, "")
        if err != nil {
                log.Debug(err)
                return err
        }
        server, err := utils.GetServerIP(serverID)
        if err != nil {
                log.Debug(err)
                return err
        }
        token, err := utils.GetServerToken(serverID)
        if err!= nil {
                log.Debug(err)
                return err
        }
        // Second join to server container
        return utils.Join(containerID, server, token, detach)
}

	

