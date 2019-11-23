package cmd

import (
	"context"

	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)
// Kill command kill a contaienr and delete all files on host

var KillCommand = cli.Command{
	Name:  "kill",
        Usage: "attach command attach io to a container or send a signal to its init process",
        ArgsUsage: `kill <-s> <container-id>`,
        Description: `Kill command kill a contaienr and delete all files on host`,
	Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "signal, s",
				Usage: `send a specific signal to container`,
			},
        },
        Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return kill(ctx, context.Args().First(),
				context.String("signal"),
			)
        },
}

func kill(ctx context.Context, containerID, signal string) error {
		log.Debug("killing a container")
		if signal == "" {
			log.Debug("default to sigterm")
			signal = "sigterm"
		}
        err := utils.KillContainer(containerID, signal)
        if err != nil {
                log.Debug(err)
                return err
		}
        return nil
}