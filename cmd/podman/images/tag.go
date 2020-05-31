package images

import (
	"github.com/containers/libpod/cmd/podman/registry"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	tagDescription = "Adds one or more additional names to locally-stored image."
	tagCommand     = &cobra.Command{
		Use:   "tag [flags] IMAGE TARGET_NAME [TARGET_NAME...]",
		Short: "Add an additional name to a local image",
		Long:  tagDescription,
		RunE:  tag,
		Args:  cobra.MinimumNArgs(2),
		Example: `podman tag 0e3bbc2 fedora:latest
  podman tag imageID:latest myNewImage:newTag
  podman tag httpd myregistryhost:5000/fedora/httpd:v2`,
	}

	imageTagCommand = &cobra.Command{
		Args:  tagCommand.Args,
		Use:   tagCommand.Use,
		Short: tagCommand.Short,
		Long:  tagCommand.Long,
		RunE:  tagCommand.RunE,
		Example: `podman image tag 0e3bbc2 fedora:latest
  podman image tag imageID:latest myNewImage:newTag
  podman image tag httpd myregistryhost:5000/fedora/httpd:v2`,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: tagCommand,
	})
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: imageTagCommand,
		Parent:  imageCmd,
	})
}

func tag(cmd *cobra.Command, args []string) error {
	return registry.ImageEngine().Tag(registry.GetContext(), args[0], args[1:], entities.ImageTagOptions{})
}
