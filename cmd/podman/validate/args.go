package validate

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NoArgs returns an error if any args are included.
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("`%s` takes no arguments", cmd.CommandPath())
	}
	return nil
}

// SubCommandExists returns an error if no sub command is provided
func SubCommandExists(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf("unrecognized command `%[1]s %[2]s`\nTry '%[1]s --help' for more information.", cmd.CommandPath(), args[0])
	}
	return errors.Errorf("missing command '%[1]s COMMAND'\nTry '%[1]s --help' for more information.", cmd.CommandPath())
}

// IdOrLatestArgs used to validate a nameOrId was provided or the "--latest" flag
func IdOrLatestArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 1 || (len(args) == 0 && !cmd.Flag("latest").Changed) {
		return fmt.Errorf("`%s` requires a name, id  or the \"--latest\" flag", cmd.CommandPath())
	}
	return nil
}
