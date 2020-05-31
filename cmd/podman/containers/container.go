package containers

import (
	"github.com/containers/libpod/cmd/podman/registry"
	"github.com/containers/libpod/cmd/podman/validate"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/libpod/pkg/util"
	"github.com/spf13/cobra"
)

var (
	// Pull in configured json library
	json = registry.JsonLibrary()

	// Command: podman _container_
	containerCmd = &cobra.Command{
		Use:              "container",
		Short:            "Manage containers",
		Long:             "Manage containers",
		TraverseChildren: true,
		RunE:             validate.SubCommandExists,
	}

	containerConfig = util.DefaultContainerConfig()
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: containerCmd,
	})
}
