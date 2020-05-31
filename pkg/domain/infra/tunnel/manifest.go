package tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containers/libpod/libpod/image"
	"github.com/containers/libpod/pkg/bindings/manifests"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/pkg/errors"
)

// ManifestCreate implements manifest create via ImageEngine
func (ir *ImageEngine) ManifestCreate(ctx context.Context, names, images []string, opts entities.ManifestCreateOptions) (string, error) {
	imageID, err := manifests.Create(ir.ClientCxt, names, images, &opts.All)
	if err != nil {
		return imageID, errors.Wrapf(err, "error creating manifest")
	}
	return imageID, err
}

// ManifestInspect returns contents of manifest list with given name
func (ir *ImageEngine) ManifestInspect(ctx context.Context, name string) ([]byte, error) {
	list, err := manifests.Inspect(ir.ClientCxt, name)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting content of manifest list or image %s", name)
	}

	buf, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		return buf, errors.Wrapf(err, "error rendering manifest for display")
	}
	return buf, err
}

// ManifestAdd adds images to the manifest list
func (ir *ImageEngine) ManifestAdd(ctx context.Context, opts entities.ManifestAddOptions) (string, error) {
	manifestAddOpts := image.ManifestAddOpts{
		All:       opts.All,
		Arch:      opts.Arch,
		Features:  opts.Features,
		Images:    opts.Images,
		OS:        opts.OS,
		OSVersion: opts.OSVersion,
		Variant:   opts.Variant,
	}
	if len(opts.Annotation) != 0 {
		annotations := make(map[string]string)
		for _, annotationSpec := range opts.Annotation {
			spec := strings.SplitN(annotationSpec, "=", 2)
			if len(spec) != 2 {
				return "", errors.Errorf("no value given for annotation %q", spec[0])
			}
			annotations[spec[0]] = spec[1]
		}
		manifestAddOpts.Annotation = annotations
	}
	listID, err := manifests.Add(ctx, opts.Images[1], manifestAddOpts)
	if err != nil {
		return listID, errors.Wrapf(err, "error adding to manifest list %s", opts.Images[1])
	}
	return listID, nil
}

// ManifestAnnotate updates an entry of the manifest list
func (ir *ImageEngine) ManifestAnnotate(ctx context.Context, names []string, opts entities.ManifestAnnotateOptions) (string, error) {
	manifestAnnotateOpts := image.ManifestAnnotateOpts{
		Arch:       opts.Arch,
		Features:   opts.Features,
		OS:         opts.OS,
		OSFeatures: opts.OSFeatures,
		OSVersion:  opts.OSVersion,
		Variant:    opts.Variant,
	}
	if len(opts.Annotation) > 0 {
		annotations := make(map[string]string)
		for _, annotationSpec := range opts.Annotation {
			spec := strings.SplitN(annotationSpec, "=", 2)
			if len(spec) != 2 {
				return "", errors.Errorf("no value given for annotation %q", spec[0])
			}
			annotations[spec[0]] = spec[1]
		}
		manifestAnnotateOpts.Annotation = annotations
	}
	updatedListID, err := manifests.Annotate(ctx, names[0], names[1], manifestAnnotateOpts)
	if err != nil {
		return updatedListID, errors.Wrapf(err, "error annotating %s of manifest list %s", names[1], names[0])
	}
	return fmt.Sprintf("%s :%s", updatedListID, names[1]), nil
}
