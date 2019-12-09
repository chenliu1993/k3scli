package cmd

import (
	"context"
	"fmt"
	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// DeployCommand deploy a pod based the config passed into server
var DeployCommand = cli.Command{
	Name:        "deploy",
	Usage:       "deploy a pod based on a yaml format file",
	ArgsUsage:   `deploy --config <yaml file> serverID`,
	Description: `Deploy command deploy a pod based the config passed into server`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config, c",
			Usage: `the pod config spec passed to server`,
		},
		&cli.BoolFlag{
			Name:  "redeploy",
			Usage: `re-deploy the pod running in cluster`,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: `force replace a pod`,
		},
	},
	Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return deploy(ctx, context.Args().First(),
			context.String("config"),
			context.Bool("redeploy"),
			context.Bool("force"),
		)
	},
}

func deploy(ctx context.Context, containerID, config string, redeploy, force bool) error {
	log.Debugf("deploying pod through %s", containerID)
	if config == "" {
		log.Debug("no config file specified")
		return fmt.Errorf("no pod config file specified")
	}
	if redeploy == true {
		err := utils.ReDeployPod(containerID, config, force)
		if err != nil {
			return err
		}
		return nil
	}
	err := utils.DeployPod(containerID, config)
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
