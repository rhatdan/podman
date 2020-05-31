package images

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/libpod/cmd/podmanV2/parse"
	"github.com/containers/libpod/cmd/podmanV2/registry"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/libpod/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	loadDescription = "Loads an image from a locally stored archive (tar file) into container storage."
	loadCommand     = &cobra.Command{
		Use:   "load [flags] [NAME[:TAG]]",
		Short: "Load an image from container archive",
		Long:  loadDescription,
		RunE:  load,
		Args:  cobra.MaximumNArgs(1),
	}
)

var (
	loadOpts entities.ImageLoadOptions
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: loadCommand,
	})

	flags := loadCommand.Flags()
	flags.StringVarP(&loadOpts.Input, "input", "i", "", "Read from specified archive file (default: stdin)")
	flags.BoolVarP(&loadOpts.Quiet, "quiet", "q", false, "Suppress the output")
	flags.StringVar(&loadOpts.SignaturePolicy, "signature-policy", "", "Pathname of signature policy file")
	if registry.IsRemote() {
		_ = flags.MarkHidden("signature-policy")
	}

}

func load(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		ref, err := reference.Parse(args[0])
		if err != nil {
			return err
		}
		if t, ok := ref.(reference.Tagged); ok {
			loadOpts.Tag = t.Tag()
		} else {
			loadOpts.Tag = "latest"
		}
		if r, ok := ref.(reference.Named); ok {
			fmt.Println(r.Name())
			loadOpts.Name = r.Name()
		}
	}
	if len(loadOpts.Input) > 0 {
		if err := parse.ValidateFileName(loadOpts.Input); err != nil {
			return err
		}
	} else {
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			return errors.Errorf("cannot read from terminal. Use command-line redirection or the --input flag.")
		}
		outFile, err := ioutil.TempFile(util.Tmpdir(), "podman")
		if err != nil {
			return errors.Errorf("error creating file %v", err)
		}
		defer os.Remove(outFile.Name())
		defer outFile.Close()

		_, err = io.Copy(outFile, os.Stdin)
		if err != nil {
			return errors.Errorf("error copying file %v", err)
		}
		loadOpts.Input = outFile.Name()
	}
	response, err := registry.ImageEngine().Load(context.Background(), loadOpts)
	if err != nil {
		return err
	}
	fmt.Println("Loaded image: " + response.Name)
	return nil
}
