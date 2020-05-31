package system

import (
	"context"
	"net/http"

	"github.com/containers/libpod/libpod/define"
	"github.com/containers/libpod/pkg/bindings"
)

// Info returns information about the libpod environment and its stores
func Info(ctx context.Context) (*define.Info, error) {
	info := define.Info{}
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(nil, http.MethodGet, "/info", nil)
	if err != nil {
		return nil, err
	}
	return &info, response.Process(&info)
}
