package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenliu1993/k3scli/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runtimeCommands = []cli.Command{
	// container-level
	cmd.RunCommand,
	cmd.JoinCommand,
	cmd.AttachCommand,
	cmd.KillCommand,
	// cluster-level
	cmd.CreateCommand,
	cmd.DeleteCommand,
	cmd.DeployCommand,
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("exiting...")
		//os.Exit(1)
	}()
	ctx := context.Background()
	cliApp := cli.NewApp()
	cliApp.Commands = runtimeCommands
	cliApp.Metadata = map[string]interface{}{
		"context": ctx,
	}
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Value: "info",
		},
	}
	cliApp.Action = func (c *cli.Context) error {
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			return err
		}
		log.SetLevel(level)
		return nil
	}
	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
