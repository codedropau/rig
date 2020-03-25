package snapshot

import (
	"context"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/docker/docker/client"
	"github.com/skpr/awsutils/ecr"

	composeconfig "github.com/codedropau/rig/internal/compose/config"
	"github.com/codedropau/rig/internal/config"
	"github.com/codedropau/rig/internal/snapshot"
)

type command struct {
	Config string

	// Authentication used when connecting to the registry.
	Username string
	Password string

	// Information used to run the correct images.
	Repository string
	Tag        string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	cfg, err := config.Load(cmd.Config)
	if err != nil {
		return err
	}

	dc, err := composeconfig.Load(cfg.Dockerfiles)
	if err != nil {
		return err
	}

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	var services []string

	for service := range cfg.Services {
		services = append(services, service)
	}

	if ecr.IsRegistry(cmd.Repository) {
		username, password, err := ecr.UpgradeAuth(cmd.Repository, cmd.Username, cmd.Password)
		if err != nil {
			return fmt.Errorf("failed to upgrade authenication for AWS ECR: %w", err)
		}

		cmd.Username = username
		cmd.Password = password
	}

	params := snapshot.Params{
		Services:   services,
		Username: cmd.Username,
		Password: cmd.Password,
		Repository: cmd.Repository,
		Tag:        cmd.Tag,
		Compose: dc,
		Config:     cfg,
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

	cmd.Flag("config", "Config file to load.").Default(".rig.yml").Envar("RIG_CONFIG").StringVar(&c.Config)

	// Information used to run the correct images.
	cmd.Flag("repository", "Tag to apply to all images when performing a snapshot.").Required().Envar("RIG_REPOSITORY").StringVar(&c.Repository)
	cmd.Flag("username", "Username used to authenticate with the registry.").Required().Envar("RIG_USERNAME").StringVar(&c.Username)
	cmd.Flag("password", "Password used to authenticate with the registry.").Required().Envar("RIG_PASSWORD").StringVar(&c.Password)
	cmd.Arg("tag", "Tag to apply to all images when performing a snapshot.").Required().StringVar(&c.Tag)
}
