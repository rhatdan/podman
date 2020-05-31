package main

import (
	"fmt"

	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/containers/libpod/pkg/adapter"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	stopCommand     cliconfig.StopValues
	stopDescription = fmt.Sprintf(`Stops one or more running containers.  The container name or ID can be used.

  A timeout to forcibly stop the container can also be set but defaults to %d seconds otherwise.`, defaultContainerConfig.Engine.StopTimeout)
	_stopCommand = &cobra.Command{
		Use:   "stop [flags] CONTAINER [CONTAINER...]",
		Short: "Stop one or more containers",
		Long:  stopDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			stopCommand.InputArgs = args
			stopCommand.GlobalFlags = MainGlobalOpts
			stopCommand.Remote = remoteclient
			return stopCmd(&stopCommand)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return checkAllLatestAndCIDFile(cmd, args, false, true)
		},
		Example: `podman stop ctrID
  podman stop --latest
  podman stop --time 2 mywebserver 6e534f14da9d`,
	}
)

func init() {
	stopCommand.Command = _stopCommand
	stopCommand.SetHelpTemplate(HelpTemplate())
	stopCommand.SetUsageTemplate(UsageTemplate())
	flags := stopCommand.Flags()
	flags.BoolVarP(&stopCommand.All, "all", "a", false, "Stop all running containers")
	flags.BoolVarP(&stopCommand.Ignore, "ignore", "i", false, "Ignore errors when a specified container is missing")
	flags.StringArrayVarP(&stopCommand.CIDFiles, "cidfile", "", nil, "Read the container ID from the file")
	flags.BoolVarP(&stopCommand.Latest, "latest", "l", false, "Act on the latest container podman is aware of")
	flags.UintVarP(&stopCommand.Timeout, "time", "t", defaultContainerConfig.Engine.StopTimeout, "Seconds to wait for stop before killing the container")
	markFlagHiddenForRemoteClient("latest", flags)
	markFlagHiddenForRemoteClient("cidfile", flags)
	markFlagHiddenForRemoteClient("ignore", flags)
	flags.SetNormalizeFunc(aliasFlags)
}

// stopCmd stops a container or containers
func stopCmd(c *cliconfig.StopValues) error {
	if c.Bool("trace") {
		span, _ := opentracing.StartSpanFromContext(Ctx, "stopCmd")
		defer span.Finish()
	}

	runtime, err := adapter.GetRuntime(getContext(), &c.PodmanCommand)
	if err != nil {
		return errors.Wrapf(err, "could not get runtime")
	}
	defer runtime.DeferredShutdown(false)

	ok, failures, err := runtime.StopContainers(getContext(), c)
	if err != nil {
		return err
	}
	return printCmdResults(ok, failures)
}
