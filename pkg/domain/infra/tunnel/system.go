package tunnel

import (
	"context"
	"errors"

	"github.com/containers/libpod/libpod/define"
	"github.com/containers/libpod/pkg/bindings/system"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/spf13/cobra"
)

func (ic *ContainerEngine) Info(ctx context.Context) (*define.Info, error) {
	return system.Info(ic.ClientCxt)
}

func (ic *ContainerEngine) RestService(_ context.Context, _ entities.ServiceOptions) error {
	panic(errors.New("rest service is not supported when tunneling"))
}

func (ic *ContainerEngine) VarlinkService(_ context.Context, _ entities.ServiceOptions) error {
	panic(errors.New("varlink service is not supported when tunneling"))
}

func (ic *ContainerEngine) SetupRootless(_ context.Context, cmd *cobra.Command) error {
	panic(errors.New("rootless engine mode is not supported when tunneling"))
}
