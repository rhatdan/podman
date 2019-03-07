package main

import (
	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	systemDescription = "Manage podman"

	systemCommand = cliconfig.PodmanCommand{
		Command: &cobra.Command{
			Use:   "system",
			Short: "Manage podman",
			Long:  systemDescription,
			RunE: func(cmd *cobra.Command, args []string) error {
				return errors.Errorf("unrecognized command `podman system %s`\nTry 'podman system --help' for more information.", args[0])
			},
		},
	}
)

var systemCommands = []*cobra.Command{
	_infoCommand,
}

func init() {
	systemCommand.AddCommand(systemCommands...)
	systemCommand.AddCommand(getSystemSubCommands()...)
	systemCommand.SetUsageTemplate(UsageTemplate())
}
