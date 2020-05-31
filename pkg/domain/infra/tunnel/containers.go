package tunnel

import (
	"context"
	"io"
	"os"

	"github.com/containers/image/v5/docker/reference"

	"github.com/containers/libpod/pkg/bindings/containers"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/pkg/errors"
)

func (ic *ContainerEngine) ContainerExists(ctx context.Context, nameOrId string) (*entities.BoolReport, error) {
	exists, err := containers.Exists(ic.ClientCxt, nameOrId)
	return &entities.BoolReport{Value: exists}, err
}

func (ic *ContainerEngine) ContainerWait(ctx context.Context, namesOrIds []string, options entities.WaitOptions) ([]entities.WaitReport, error) {
	var (
		responses []entities.WaitReport
	)
	cons, err := getContainersByContext(ic.ClientCxt, false, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, c := range cons {
		response := entities.WaitReport{Id: c.ID}
		exitCode, err := containers.Wait(ic.ClientCxt, c.ID, &options.Condition)
		if err != nil {
			response.Error = err
		} else {
			response.ExitCode = exitCode
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (ic *ContainerEngine) ContainerPause(ctx context.Context, namesOrIds []string, options entities.PauseUnPauseOptions) ([]*entities.PauseUnpauseReport, error) {
	var (
		reports []*entities.PauseUnpauseReport
	)
	ctrs, err := getContainersByContext(ic.ClientCxt, options.All, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, c := range ctrs {
		err := containers.Pause(ic.ClientCxt, c.ID)
		reports = append(reports, &entities.PauseUnpauseReport{Id: c.ID, Err: err})
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerUnpause(ctx context.Context, namesOrIds []string, options entities.PauseUnPauseOptions) ([]*entities.PauseUnpauseReport, error) {
	var (
		reports []*entities.PauseUnpauseReport
	)
	ctrs, err := getContainersByContext(ic.ClientCxt, options.All, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, c := range ctrs {
		err := containers.Unpause(ic.ClientCxt, c.ID)
		reports = append(reports, &entities.PauseUnpauseReport{Id: c.ID, Err: err})
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerStop(ctx context.Context, namesOrIds []string, options entities.StopOptions) ([]*entities.StopReport, error) {
	var (
		reports []*entities.StopReport
	)
	ctrs, err := getContainersByContext(ic.ClientCxt, options.All, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, c := range ctrs {
		report := entities.StopReport{Id: c.ID}
		report.Err = containers.Stop(ic.ClientCxt, c.ID, &options.Timeout)
		// TODO we need to associate errors returned by http with common
		// define.errors so that we can equity tests. this will allow output
		// to be the same as the native client
		reports = append(reports, &report)
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerKill(ctx context.Context, namesOrIds []string, options entities.KillOptions) ([]*entities.KillReport, error) {
	var (
		reports []*entities.KillReport
	)
	ctrs, err := getContainersByContext(ic.ClientCxt, options.All, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, c := range ctrs {
		reports = append(reports, &entities.KillReport{
			Id:  c.ID,
			Err: containers.Kill(ic.ClientCxt, c.ID, options.Signal),
		})
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerRestart(ctx context.Context, namesOrIds []string, options entities.RestartOptions) ([]*entities.RestartReport, error) {
	var (
		reports []*entities.RestartReport
		timeout *int
	)
	if options.Timeout != nil {
		t := int(*options.Timeout)
		timeout = &t
	}
	ctrs, err := getContainersByContext(ic.ClientCxt, options.All, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, c := range ctrs {
		reports = append(reports, &entities.RestartReport{
			Id:  c.ID,
			Err: containers.Restart(ic.ClientCxt, c.ID, timeout),
		})
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerRm(ctx context.Context, namesOrIds []string, options entities.RmOptions) ([]*entities.RmReport, error) {
	var (
		reports []*entities.RmReport
	)
	ctrs, err := getContainersByContext(ic.ClientCxt, options.All, namesOrIds)
	if err != nil {
		return nil, err
	}
	// TODO there is no endpoint for container eviction.  Need to discuss
	for _, c := range ctrs {
		reports = append(reports, &entities.RmReport{
			Id:  c.ID,
			Err: containers.Remove(ic.ClientCxt, c.ID, &options.Force, &options.Volumes),
		})
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerInspect(ctx context.Context, namesOrIds []string, options entities.InspectOptions) ([]*entities.ContainerInspectReport, error) {
	var (
		reports []*entities.ContainerInspectReport
	)
	ctrs, err := getContainersByContext(ic.ClientCxt, false, namesOrIds)
	if err != nil {
		return nil, err
	}
	for _, con := range ctrs {
		data, err := containers.Inspect(ic.ClientCxt, con.ID, &options.Size)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &entities.ContainerInspectReport{InspectContainerData: data})
	}
	return reports, nil
}

func (ic *ContainerEngine) ContainerTop(ctx context.Context, options entities.TopOptions) (*entities.StringSliceReport, error) {
	switch {
	case options.Latest:
		return nil, errors.New("latest is not supported")
	case options.NameOrID == "":
		return nil, errors.New("NameOrID must be specified")
	}

	topOutput, err := containers.Top(ic.ClientCxt, options.NameOrID, options.Descriptors)
	if err != nil {
		return nil, err
	}
	return &entities.StringSliceReport{Value: topOutput}, nil
}

func (ic *ContainerEngine) ContainerCommit(ctx context.Context, nameOrId string, options entities.CommitOptions) (*entities.CommitReport, error) {
	var (
		repo string
		tag  string = "latest"
	)
	if len(options.ImageName) > 0 {
		ref, err := reference.Parse(options.ImageName)
		if err != nil {
			return nil, err
		}
		if t, ok := ref.(reference.Tagged); ok {
			tag = t.Tag()
		}
		if r, ok := ref.(reference.Named); ok {
			repo = r.Name()
		}
		if len(repo) < 1 {
			return nil, errors.Errorf("invalid image name %q", options.ImageName)
		}
	}
	commitOpts := containers.CommitOptions{
		Author:  &options.Author,
		Changes: options.Changes,
		Comment: &options.Message,
		Format:  &options.Format,
		Pause:   &options.Pause,
		Repo:    &repo,
		Tag:     &tag,
	}
	response, err := containers.Commit(ic.ClientCxt, nameOrId, commitOpts)
	if err != nil {
		return nil, err
	}
	return &entities.CommitReport{Id: response.ID}, nil
}

func (ic *ContainerEngine) ContainerExport(ctx context.Context, nameOrId string, options entities.ContainerExportOptions) error {
	var (
		err error
		w   io.Writer
	)
	if len(options.Output) > 0 {
		w, err = os.Create(options.Output)
		if err != nil {
			return err
		}
	}
	return containers.Export(ic.ClientCxt, nameOrId, w)
}
