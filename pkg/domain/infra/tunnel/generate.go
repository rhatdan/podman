package tunnel

import (
	"context"

	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/pkg/errors"
)

func (ic *ContainerEngine) GenerateSystemd(ctx context.Context, nameOrID string, options entities.GenerateSystemdOptions) (*entities.GenerateSystemdReport, error) {
	return nil, errors.New("not implemented for tunnel")
}
