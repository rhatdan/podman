package main

import (
	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	podDescription = `Pods are a group of one or more containers sharing the same network, pid and ipc namespaces.`
)
var podCommand = cliconfig.PodmanCommand{
	Command: &cobra.Command{
		Use:   "pod",
		Short: "Manage pods",
		Long:  podDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Errorf("unrecognized command `podman pod %s`\nTry 'podman pod --help' for more information.", args[0])
		},
	},
}

//podSubCommands are implemented both in local and remote clients
var podSubCommands = []*cobra.Command{
	_podCreateCommand,
	_podExistsCommand,
	_podInspectCommand,
	_podKillCommand,
	_podPauseCommand,
	_podPsCommand,
	_podRestartCommand,
	_podRmCommand,
	_podStartCommand,
	_podStatsCommand,
	_podStopCommand,
	_podTopCommand,
	_podUnpauseCommand,
}

func init() {
	podCommand.AddCommand(podSubCommands...)
	podCommand.SetHelpTemplate(HelpTemplate())
	podCommand.SetUsageTemplate(UsageTemplate())
}
