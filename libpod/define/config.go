package define

var (
	// DefaultInfraImage to use for infra container
	DefaultInfraImage = "k8s.gcr.io/pause:3.2"
	// DefaultInfraCommand to be run in an infra container
	DefaultInfraCommand = "/pause"
	// DefaultSHMLockPath is the default path for SHM locks
	DefaultSHMLockPath = "/libpod_lock"
	// DefaultRootlessSHMLockPath is the default path for rootless SHM locks
	DefaultRootlessSHMLockPath = "/libpod_rootless_lock"
)

const (
	// DefaultTransport is a prefix that we apply to an image name
	// to check docker hub first for the image
	DefaultTransport = "docker://"
)

// InfoData holds the info type, i.e store, host etc and the data for each type
type InfoData struct {
	Type string
	Data map[string]interface{}
}

// VolumeDriverLocal is the "local" volume driver. It is managed by libpod
// itself.
const VolumeDriverLocal = "local"

const (
	OCIManifestDir  = "oci-dir"
	OCIArchive      = "oci-archive"
	V2s2ManifestDir = "docker-dir"
	V2s2Archive     = "docker-archive"
)
