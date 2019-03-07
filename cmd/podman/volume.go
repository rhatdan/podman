package main

import (
	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var volumeDescription = `Volumes are created in and can be shared between containers.`

var volumeCommand = cliconfig.PodmanCommand{
	Command: &cobra.Command{
		Use:   "volume",
		Short: "Manage volumes",
		Long:  volumeDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Errorf("unrecognized command `podman volume %s`\nTry 'podman volume --help' for more information.", args[0])
		},
	},
}
var volumeSubcommands = []*cobra.Command{
	_volumeCreateCommand,
	_volumeLsCommand,
	_volumeRmCommand,
	_volumeInspectCommand,
	_volumePruneCommand,
}

func init() {
	volumeCommand.SetUsageTemplate(UsageTemplate())
	volumeCommand.AddCommand(volumeSubcommands...)
	rootCmd.AddCommand(volumeCommand.Command)
}
