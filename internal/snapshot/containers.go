package snapshot

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	dockerauth "github.com/codedropau/rig/internal/docker/auth"
	"github.com/codedropau/rig/internal/docker/stream"
)

// Helper function to snapshot containers associated with a project.
func snapshotContainers(ctx context.Context, cli *client.Client, params Params) error {
	projectFilter := filters.NewArgs()
	projectFilter.Add("label", fmt.Sprintf("%s=%s", LabelProject, params.Config.Project))

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: projectFilter,
	})
	if err != nil {
		return err
	}

	for _, service := range params.Services {
		c, err := getContainer(containers, service)
		if err != nil {
			return err
		}

		reference := fmt.Sprintf("%s:%s-service-%s", params.Repository, params.Tag, c.Labels[LabelService])

		fmt.Printf("Snapshotting container '%s' to '%s'\n", c.Labels[LabelService], reference)

		_, err = cli.ContainerCommit(ctx, c.ID, types.ContainerCommitOptions{
			Reference: reference,
			Pause:     false,
		})
		if err != nil {
			return err
		}

		auth, err := dockerauth.Base64(params.Username, params.Password)
		if err != nil {
			return err
		}

		push, err := cli.ImagePush(ctx, reference, types.ImagePushOptions{
			RegistryAuth: auth,
		})
		if err != nil {
			return err
		}

		err = stream.Push(os.Stdout, push)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to return a container if it exists, or error if not found. This ensures a service exists when
// performing a snapshot.
func getContainer(containers []types.Container, service string) (types.Container, error) {
	for _, container := range containers {
		if _, ok := container.Labels[LabelService]; !ok {
			continue
		}

		if container.Labels[LabelService] != service {
			continue
		}

		return container, nil
	}

	return types.Container{}, fmt.Errorf("not found: %s", service)
}
