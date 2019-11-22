package snapshot

import (
	"context"

	"github.com/alecthomas/kingpin"
	"github.com/docker/docker/client"

	"github.com/nickschuch/rig/internal/compose"
	"github.com/nickschuch/rig/internal/snapshot"
	"github.com/nickschuch/rig/internal/config"
)

type command struct {
	Project string
	Repository string
	Tag string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	cfg := config.Config{
		Services: []string{
			"nginx",
			"php-fpm",
			"mysql-default",
		},
	}

	project, err := compose.Project(cmd.Project)
	if err != nil {
		return err
	}

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	params := snapshot.Params {
		Project: project,
		Services: cfg.Services,
		Repository: cmd.Repository,
		Tag: cmd.Tag,
	}

	err = snapshot.All(ctx, cli, params)
	if err != nil {
		return err
	}

	return nil
}

// Command which snapshots a Docker Compose stack.
func Command(app *kingpin.Application) {
	c := new(command)

	cmd := app.Command("snapshot", "Takes a snapshot of the existing Docker Compose stack.").Action(c.run)
	cmd.Flag("project", "Tag to apply to all images when performing a snapshot.").Envar("RIG_PROJECT").StringVar(&c.Project)
	cmd.Flag("repository", "Tag to apply to all images when performing a snapshot.").Required().Envar("RIG_REPOSITORY").StringVar(&c.Repository)
	cmd.Arg("tag", "Tag to apply to all images when performing a snapshot.").Required().StringVar(&c.Tag)
}