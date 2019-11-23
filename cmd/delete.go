package cmd


//delete command deletes a cluster and rm all files on hosts

import (
	"context"
	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

)
// DeleteCommand deletes a cluster
var DeleteCommand = cli.Command{
	Name:  "delete",
        Usage: "delete a cluster based on cluster name",
        ArgsUsage: `delete <cluster-name>`,
        Description: `delete a cluster named <cluster-name>`,
	Flags: []cli.Flag{
        },
        Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return deletecluster(ctx, context.Args().First())
        },
}

func deletecluster(ctx context.Context, clusterName string) error {
		log.Debug("Begin deleting a cluster")
        err := utils.DeleteCluster(clusterName)
        if err != nil {
                log.Debug(err)
                return err
        }
        return nil
}