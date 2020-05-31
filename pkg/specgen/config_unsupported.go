// +build !linux

package specgen

import (
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

func (s *SpecGenerator) getSeccompConfig(configSpec *spec.Spec) (*spec.LinuxSeccomp, error) {
	return nil, errors.New("function not supported on non-linux OS's")
}
