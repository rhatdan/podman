package main

import (
	"fmt"

	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/containers/libpod/libpod/define"
	"github.com/containers/libpod/pkg/adapter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rmCommand     cliconfig.RmValues
	rmDescription = fmt.Sprintf(`Removes one or more containers from the host. The container name or ID can be used.

  Command does not remove images. Running or unusable containers will not be removed without the -f option.`)
	_rmCommand = &cobra.Command{
		Use:   "rm [flags] CONTAINER [CONTAINER...]",
		Short: "Remove one or more containers",
		Long:  rmDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			rmCommand.InputArgs = args
			rmCommand.GlobalFlags = MainGlobalOpts
			rmCommand.Remote = remoteclient
			return rmCmd(&rmCommand)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return checkAllLatestAndCIDFile(cmd, args, false, true)
		},
		Example: `podman rm imageID
  podman rm mywebserver myflaskserver 860a4b23
  podman rm --force --all
  podman rm -f c684f0d469f2`,
	}
)

func init() {
	rmCommand.Command = _rmCommand
	rmCommand.SetHelpTemplate(HelpTemplate())
	rmCommand.SetUsageTemplate(UsageTemplate())
	flags := rmCommand.Flags()
	flags.BoolVarP(&rmCommand.All, "all", "a", false, "Remove all containers")
	flags.BoolVarP(&rmCommand.Ignore, "ignore", "i", false, "Ignore errors when a specified container is missing")
	flags.BoolVarP(&rmCommand.Force, "force", "f", false, "Force removal of a running or unusable container.  The default is false")
	flags.BoolVarP(&rmCommand.Latest, "latest", "l", false, "Act on the latest container podman is aware of")
	flags.BoolVar(&rmCommand.Storage, "storage", false, "Remove container from storage library")
	flags.BoolVarP(&rmCommand.Volumes, "volumes", "v", false, "Remove anonymous volumes associated with the container")
	flags.StringArrayVarP(&rmCommand.CIDFiles, "cidfile", "", nil, "Read the container ID from the file")
	markFlagHiddenForRemoteClient("ignore", flags)
	markFlagHiddenForRemoteClient("cidfile", flags)
	markFlagHiddenForRemoteClient("latest", flags)
	markFlagHiddenForRemoteClient("storage", flags)
}

// rmCmd removes one or more containers
func rmCmd(c *cliconfig.RmValues) error {
	runtime, err := adapter.GetRuntime(getContext(), &c.PodmanCommand)
	if err != nil {
		return errors.Wrapf(err, "could not get runtime")
	}
	defer runtime.DeferredShutdown(false)

	// Storage conflicts with --all/--latest/--volumes/--cidfile/--ignore
	if c.Storage {
		if c.All || c.Ignore || c.Latest || c.Volumes || c.CIDFiles != nil {
			return errors.Errorf("--storage conflicts with --volumes, --all, --latest, --ignore and --cidfile")
		}
	}

	ok, failures, err := runtime.RemoveContainers(getContext(), c)
	if err != nil {
		if len(c.InputArgs) < 2 {
			exitCode = setExitCode(err)
		}
		return err
	}

	if len(failures) > 0 {
		for _, err := range failures {
			if errors.Cause(err) == define.ErrWillDeadlock {
				logrus.Errorf("Potential deadlock detected - please run 'podman system renumber' to resolve")
			}
			exitCode = setExitCode(err)
		}
	}

	return printCmdResults(ok, failures)
}
