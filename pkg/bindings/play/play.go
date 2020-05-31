package play

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/containers/image/v5/types"
	"github.com/containers/libpod/pkg/bindings"
	"github.com/containers/libpod/pkg/domain/entities"
)

func PlayKube(ctx context.Context, path string, options entities.PlayKubeOptions) (*entities.PlayKubeReport, error) {
	var report entities.PlayKubeReport
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	params := url.Values{}
	params.Set("network", options.Network)
	if options.SkipTLSVerify != types.OptionalBoolUndefined {
		params.Set("tlsVerify", strconv.FormatBool(options.SkipTLSVerify == types.OptionalBoolTrue))
	}

	response, err := conn.DoRequest(f, http.MethodPost, "/play/kube", params)
	if err != nil {
		return nil, err
	}
	if err := response.Process(&report); err != nil {
		return nil, err
	}

	return &report, nil
}
