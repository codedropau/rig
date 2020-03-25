package snapshot

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	composeconfig "github.com/codedropau/rig/internal/compose/config"
	stringutils "github.com/codedropau/rig/internal/utils/string"
	dockerauth "github.com/codedropau/rig/internal/docker/auth"
	"github.com/codedropau/rig/internal/docker/stream"
	"github.com/codedropau/rig/internal/utils/file"
)

// Helper function to snapshot volumes associated with a project.
func snapshotVolumes(ctx context.Context, cli *client.Client, params Params) error {
	projectFilter := filters.NewArgs()
	projectFilter.Add("label", fmt.Sprintf("%s=%s", LabelProject, params.Config.Project))

	list, err := cli.VolumeList(ctx, projectFilter)
	if err != nil {
		return err
	}

	for _, volume := range list.Volumes {
		if _, ok := volume.Labels[LabelVolume]; !ok {
			continue
		}

		name := volume.Labels[LabelVolume]

		if !referencedInServices(params.Compose.Services, params.Services, name) {
			continue
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		dockerfilePath := fmt.Sprintf("%s/.rig/volume/%s.dockerfile", dir, name)

		tag := fmt.Sprintf("%s-volume-%s", params.Tag, name)

		if _, ok := volume.Options["device"]; !ok {
			return fmt.Errorf("cannot find device for volume: %s", name)
		}

		fmt.Printf("Snapshotting volume '%s' to '%s:%s'\n", name, params.Repository, tag)

		tmpl := "FROM %s\nWORKDIR /volume\nADD --chown=%s:%s %s /volume\nRUN chmod -R 775 /volume"

		err = file.Write(dockerfilePath, fmt.Sprintf(tmpl, params.Config.Volume.From, params.Config.Volume.User, params.Config.Volume.Group, volume.Options["device"]))
		if err != nil {
			return err
		}
		defer os.Remove(dockerfilePath)

		options := &archive.TarOptions{
			IncludeFiles: []string{
				dockerfilePath,
				volume.Options["device"],
			},
		}

		reader, err := archive.TarWithOptions("/", options)
		if err != nil {
			return err
		}

		reference := fmt.Sprintf("%s:%s", params.Repository, tag)

		output, err := cli.ImageBuild(ctx, reader, types.ImageBuildOptions{
			Dockerfile: dockerfilePath,
			Tags: []string{
				reference,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to build image: %w", err)
		}

		err = stream.Build(os.Stdout, output.Body)
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

// Helper function to check if a volume is referenced in the list of Services requested by the user.
func referencedInServices(services map[string]*composeconfig.Service, allowed []string, volume string) bool {
	for name, service := range services {
		if !stringutils.Contains(allowed, name) {
			continue
		}

		if service.MountsVolume(volume) {
			return true
		}
	}

	return false
}