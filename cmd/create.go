package cmd

// Create command does not create  a command, it is used
// to create a cluster
import (
	"context"
	clusterconfig "github.com/chenliu1993/k3scli/pkg/config/cluster"
	"github.com/chenliu1993/k3scli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"sort"
)

// CreateCommand creates a cluster
var CreateCommand = cli.Command{
	Name:        "create",
	Usage:       "create a cluster based on config file",
	ArgsUsage:   `create <cluster-name>`,
	Description: `create a cluster named <cluster-name>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: `the config file used to created a cluster`,
		},
		&cli.StringFlag{
			Name:  "mode",
			Usage: `containerd or docker`,
		},
	},
	Action: func(context *cli.Context) error {
		ctx, err := cliContextToContext(context)
		if err != nil {
			return err
		}
		return createcluster(ctx, context.Args().First(),
			context.String("config"),
			context.String("mode"),
		)
	},
}

func createcluster(ctx context.Context, clusterName, config, mode string) error {
	log.Debug("Begin creating a cluster")
	var cluster clusterconfig.Cluster
	var err error
	if config == "" {
		log.Debug("no config is specified, default config will be used")
		cluster = clusterconfig.DefaultClusterConfig
	} else {
		cluster, err = clusterconfig.Load(config)
		if err != nil {
			log.Fatal(err)
		}
	}
	sort.Slice(cluster.Nodes, func(i int, j int) bool {
		return cluster.Nodes[i].Label < cluster.Nodes[j].Label
	})
	err = utils.CreateCluster(clusterName, mode, cluster)
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
