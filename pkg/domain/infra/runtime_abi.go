// +build ABISupport

package infra

import (
	"context"
	"fmt"

	"github.com/containers/libpod/libpod"
	"github.com/containers/libpod/pkg/bindings"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/libpod/pkg/domain/infra/abi"
	"github.com/containers/libpod/pkg/domain/infra/tunnel"
)

// NewContainerEngine factory provides a libpod runtime for container-related operations
func NewContainerEngine(facts *entities.PodmanConfig) (entities.ContainerEngine, error) {
	switch facts.EngineMode {
	case entities.ABIMode:
		r, err := NewLibpodRuntime(facts.FlagSet, facts)
		return r, err
	case entities.TunnelMode:
		ctx, err := bindings.NewConnection(context.Background(), facts.Uri, facts.Identities...)
		return &tunnel.ContainerEngine{ClientCxt: ctx}, err
	}
	return nil, fmt.Errorf("runtime mode '%v' is not supported", facts.EngineMode)
}

// NewContainerEngine factory provides a libpod runtime for image-related operations
func NewImageEngine(facts *entities.PodmanConfig) (entities.ImageEngine, error) {
	switch facts.EngineMode {
	case entities.ABIMode:
		r, err := NewLibpodImageRuntime(facts.FlagSet, facts)
		return r, err
	case entities.TunnelMode:
		ctx, err := bindings.NewConnection(context.Background(), facts.Uri, facts.Identities...)
		return &tunnel.ImageEngine{ClientCxt: ctx}, err
	}
	return nil, fmt.Errorf("runtime mode '%v' is not supported", facts.EngineMode)
}

// NewSystemEngine factory provides a libpod runtime for specialized system operations
func NewSystemEngine(setup entities.EngineSetup, facts *entities.PodmanConfig) (entities.SystemEngine, error) {
	switch facts.EngineMode {
	case entities.ABIMode:
		var r *libpod.Runtime
		var err error
		switch setup {
		case entities.NormalMode:
			r, err = GetRuntime(context.Background(), facts.FlagSet, facts)
		case entities.RenumberMode:
			r, err = GetRuntimeRenumber(context.Background(), facts.FlagSet, facts)
		case entities.ResetMode:
			r, err = GetRuntimeRenumber(context.Background(), facts.FlagSet, facts)
		case entities.MigrateMode:
			name, flagErr := facts.FlagSet.GetString("new-runtime")
			if flagErr != nil {
				return nil, flagErr
			}
			r, err = GetRuntimeMigrate(context.Background(), facts.FlagSet, facts, name)
		case entities.NoFDsMode:
			r, err = GetRuntimeDisableFDs(context.Background(), facts.FlagSet, facts)
		}
		return &abi.SystemEngine{Libpod: r}, err
	case entities.TunnelMode:
		return nil, fmt.Errorf("tunnel system runtime not supported")
	}
	return nil, fmt.Errorf("runtime mode '%v' is not supported", facts.EngineMode)
}
