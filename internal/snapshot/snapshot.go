package snapshot

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"

	composeconfig "github.com/codedropau/rig/internal/compose/config"
	"github.com/codedropau/rig/internal/config"
)

// Params passed to the All function.
type Params struct {
	Services   []string
	Username string
	Password string
	Repository string
	Tag        string
	Compose *composeconfig.Config
	Config     *config.Config
}

// All containers and volumes related to a project.
func All(ctx context.Context, cli *client.Client, params Params) error {
	err := snapshotContainers(ctx, cli, params)
	if err != nil {
		return fmt.Errorf("failed to snapshot container: %w", err)
	}

	err = snapshotVolumes(ctx, cli, params)
	if err != nil {
		return fmt.Errorf("failed to snapshot volume: %w", err)
	}

	return nil
}
