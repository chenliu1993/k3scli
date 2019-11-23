package cmd

import (
	"context"

	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)
// Attach command attach io to a container

var AttachCommand = cli.Command{
	Name:  "attach",
        Usage: "attach command attach io to a container",
        ArgsUsage: `attach <container-id>`,
        Description: `attach to a running container and set stdio`,
	Flags: []cli.Flag{
        },
        Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return attach(ctx, context.Args().First())
        },
}

func attach(ctx context.Context, containerID string) error {
        log.Debug("Begin attach to a existing container, first checking args")
        err := utils.AttachContainer(containerID)
        if err != nil {
                log.Debug(err)
                return err
        }
        return nil
}