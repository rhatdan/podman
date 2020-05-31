package compat

import (
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/containers/storage/pkg/archive"
)

// Create container
// swagger:response ContainerCreateResponse
type swagCtrCreateResponse struct {
	// in:body
	Body struct {
		entities.ContainerCreateResponse
	}
}

// Wait container
// swagger:response ContainerWaitResponse
type swagCtrWaitResponse struct {
	// in:body
	Body struct {
		// container exit code
		StatusCode int
		Error      struct {
			Message string
		}
	}
}

// Object Changes
// swagger:response Changes
type swagChangesResponse struct {
	// in:body
	Body struct {
		Changes []archive.Change
	}
}
