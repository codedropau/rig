package snapshot

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type Params struct {
	Project string
	Services []string
	Repository string
	Tag string
}

func All(ctx context.Context, cli *client.Client, params Params) error {
	err := snapshotContainers(ctx, cli, params)
	if err != nil {
		return errors.Wrap(err, "failed to snapshot container")
	}

	err = Volumes(ctx, cli, params)
	if err != nil {
		return errors.Wrap(err, "failed to snapshot volume")
	}

	return nil
}
