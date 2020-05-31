package pods

import (
	"github.com/containers/libpod/cmd/podman/registry"
	"github.com/containers/libpod/cmd/podman/validate"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/libpod/pkg/util"
	"github.com/spf13/cobra"
)

var (
	// Command: podman _generate_
	generateCmd = &cobra.Command{
		Use:              "generate",
		Short:            "Generate structured data based on containers and pods.",
		Long:             "Generate structured data (e.g., Kubernetes yaml or systemd units) based on containers and pods.",
		TraverseChildren: true,
		RunE:             validate.SubCommandExists,
	}
	containerConfig = util.DefaultContainerConfig()
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode},
		Command: generateCmd,
	})
}
