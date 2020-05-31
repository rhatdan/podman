package manifest

import (
	"github.com/containers/libpod/cmd/podman/registry"
	"github.com/containers/libpod/cmd/podman/validate"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	manifestDescription = "Creates, modifies, and pushes manifest lists and image indexes."
	manifestCmd         = &cobra.Command{
		Use:              "manifest",
		Short:            "Manipulate manifest lists and image indexes",
		Long:             manifestDescription,
		TraverseChildren: true,
		RunE:             validate.SubCommandExists,
		Example: `podman manifest add mylist:v1.11 image:v1.11-amd64
  podman manifest create localhost/list
  podman manifest inspect localhost/list
  podman manifest annotate --annotation left=right mylist:v1.11 image:v1.11-amd64`,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: manifestCmd,
	})
}
