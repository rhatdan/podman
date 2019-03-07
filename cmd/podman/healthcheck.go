package main

import (
	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var healthcheckDescription = "Manage health checks on containers"
var healthcheckCommand = cliconfig.PodmanCommand{
	Command: &cobra.Command{
		Use:   "healthcheck",
		Short: "Manage Healthcheck",
		Long:  healthcheckDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Errorf("unrecognized command `podman healthcheck %s`\nTry 'podman healthcheck --help' for more information.", args[0])
		},
	},
}

// Commands that are universally implemented
var healthcheckCommands []*cobra.Command

func init() {
	healthcheckCommand.AddCommand(healthcheckCommands...)
	healthcheckCommand.AddCommand(getHealtcheckSubCommands()...)
	healthcheckCommand.SetUsageTemplate(UsageTemplate())
	rootCmd.AddCommand(healthcheckCommand.Command)
}
