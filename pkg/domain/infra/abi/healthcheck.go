package abi

import (
	"context"

	"github.com/containers/libpod/libpod"
	"github.com/containers/libpod/libpod/define"
	"github.com/containers/libpod/pkg/domain/entities"
)

func (ic *ContainerEngine) HealthCheckRun(ctx context.Context, nameOrId string, options entities.HealthCheckOptions) (*define.HealthCheckResults, error) {
	status, err := ic.Libpod.HealthCheck(nameOrId)
	if err != nil {
		return nil, err
	}
	hcStatus := "unhealthy"
	if status == libpod.HealthCheckSuccess {
		hcStatus = "healthy"
	}
	report := define.HealthCheckResults{
		Status: hcStatus,
	}
	return &report, nil
}
