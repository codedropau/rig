package snapshot

import (
	"context"
	"fmt"
	"github.com/codedropau/rig/internal/config"

	"github.com/docker/docker/client"
)

// Params passed to the All function.
type Params struct {
	Services   []string
	Repository string
	Tag        string
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
