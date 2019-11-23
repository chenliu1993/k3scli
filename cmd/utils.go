package cmd

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func cliContextToContext(c *cli.Context) (context.Context, error) {
	log.Debug("copy ctx")
	if c == nil {
		return nil, errors.New("need cli.Context")
	}

	// extract the main context
	ctx, ok := c.App.Metadata["context"].(context.Context)
	if !ok {
		return nil, errors.New("invalid or missing context in metadata")
	}

	return ctx, nil
}
