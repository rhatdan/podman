package main

import (
	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	generateCommand     cliconfig.PodmanCommand
	generateDescription = "Generate structured data based for a containers and pods"
	_generateCommand    = &cobra.Command{
		Use:   "generate",
		Short: "Generated structured data",
		Long:  generateDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Errorf("unrecognized command `podman generate %s`\nTry 'podman generate --help' for more information.", args[0])
		},
	}
)

func init() {
	generateCommand.Command = _generateCommand
	generateCommand.AddCommand(getGenerateSubCommands()...)
	generateCommand.SetUsageTemplate(UsageTemplate())
}
