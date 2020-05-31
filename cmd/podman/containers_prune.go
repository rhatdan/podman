package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/containers/libpod/cmd/podman/cliconfig"
	"github.com/containers/libpod/cmd/podman/shared"
	"github.com/containers/libpod/libpod/define"
	"github.com/containers/libpod/pkg/adapter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	pruneContainersCommand     cliconfig.PruneContainersValues
	pruneContainersDescription = `
	podman container prune

	Removes all stopped | exited containers
`
	_pruneContainersCommand = &cobra.Command{
		Use:   "prune",
		Args:  noSubArgs,
		Short: "Remove all stopped | exited containers",
		Long:  pruneContainersDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			pruneContainersCommand.InputArgs = args
			pruneContainersCommand.GlobalFlags = MainGlobalOpts
			pruneContainersCommand.Remote = remoteclient
			return pruneContainersCmd(&pruneContainersCommand)
		},
	}
)

func init() {
	pruneContainersCommand.Command = _pruneContainersCommand
	pruneContainersCommand.SetHelpTemplate(HelpTemplate())
	pruneContainersCommand.SetUsageTemplate(UsageTemplate())
	flags := pruneContainersCommand.Flags()
	flags.BoolVarP(&pruneContainersCommand.Force, "force", "f", false, "Skip interactive prompt for container removal")
	flags.StringArrayVar(&pruneContainersCommand.Filter, "filter", []string{}, "Provide filter values (e.g. 'until=<timestamp>')")
}

func pruneContainersCmd(c *cliconfig.PruneContainersValues) error {
	if !c.Force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf(`WARNING! This will remove all stopped containers.
Are you sure you want to continue? [y/N] `)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return errors.Wrapf(err, "error reading input")
		}
		if strings.ToLower(answer)[0] != 'y' {
			return nil
		}
	}

	runtime, err := adapter.GetRuntime(getContext(), &c.PodmanCommand)
	if err != nil {
		return errors.Wrapf(err, "could not get runtime")
	}
	defer runtime.DeferredShutdown(false)

	maxWorkers := shared.DefaultPoolSize("prune")
	if c.GlobalIsSet("max-workers") {
		maxWorkers = c.GlobalFlags.MaxWorks
	}
	ok, failures, err := runtime.Prune(getContext(), maxWorkers, c.Filter)
	if err != nil {
		if errors.Cause(err) == define.ErrNoSuchCtr {
			if len(c.InputArgs) > 1 {
				exitCode = define.ExecErrorCodeGeneric
			} else {
				exitCode = 1
			}
		}
		return err
	}
	if len(failures) > 0 {
		exitCode = define.ExecErrorCodeGeneric
	}
	return printCmdResults(ok, failures)
}
