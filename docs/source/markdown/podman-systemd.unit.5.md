% podman-systemd.unit 5

## NAME

podman\-systemd.unit - systemd units using Podman

## SYNOPSIS

*name*.container, *name*.volume

### Podman unit search path

 * /etc/containers/systemd/
 * /usr/share/containers/systemd/

### Podman user unit search path

 * $XDG_CONFIG_HOME/containers/systemd/
 * ~/.config/containers/systemd/

## DESCRIPTION

Podman supports starting containers (and creating volumes) via systemd by using a
[systemd generator](https://www.freedesktop.org/software/systemd/man/systemd.generator.html).
These files are read during boot (and when `systemctl daemon-reload` is run) and generate
corresponding regular systemd service unit files. Both system and user systemd units are supported.

The Podman generator reads the search paths above and reads files with the extensions `.container`
and `.volume`, and for each file generates a similarly named `.service` file. These units can be started
and managed with systemctl like any other systemd service.

The Podman files use the same format as [regular systemd unit files](https://www.freedesktop.org/software/systemd/man/systemd.syntax.html).
Each file type has a custom section (for example, `[Container]`) that is handled by Podman, and all
other sections will be passed on untouched, allowing the use of any normal systemd configuration options
like dependencies or cgroup limits.

### Enabling unit files

The services created by Podman are considered transient by systemd, which means they don't have the same
persistence rules as regular units. In particular, it is not possible to "systemctl enable" them
in order for them to become automatically enabled on the next boot.

To compensate for this, the generator manually applies the `[Install]` section of the container definition
unit files during generation, in the same way `systemctl enable` would do when run later.

For example, to start a container on boot, add something like this to the file:

```
[Install]
WantedBy=multi-user.target
```

Currently, only the `Alias`, `WantedBy` and `RequiredBy` keys are supported.

**NOTE:** To express dependencies between containers, use the generated names of the service. In other
words `WantedBy=other.service`, not `WantedBy=other.container`. The same is
true for other kinds of dependencies, too, like `After=other.service`.

### Container units

Container units are named with a `.container` extension and contain a `[Container] `section describing
the container that should be run as a service. The resulting service file will contain a line like
`ExecStart=podman run … image-name`, and most of the keys in this section control the command-line
options passed to Podman. However, some options also affect the details of how systemd is set up to run and
interact with the container.

By default, the Podman container will have the same name as the unit, but with a `systemd-` prefix.
I.e. a `$name.container` file will create a `$name.service` unit and a `systemd-$name` Podman container.

There is only one required key, `Image`, which defines the container image the service should run.

Supported keys in `Container` section are:

#### `Image=`

The image to run in the container. This image must be locally installed for the service to work
when it is activated, because the generated service file will never try to download images.
It is recommended to use a fully qualified image name rather than a short name, both for
performance and robustness reasons.

The format of the name is the same as when passed to `podman run`, so it supports e.g., using
`:tag` or using digests guarantee a specific image version.

#### `ContainerName=`

The (optional) name of the Podman container. If this is not specified, the default value
of `systemd-%N` will be used, which is the same as the service name but with a `systemd-`
prefix to avoid conflicts with user-managed containers.

#### `Environment=`

Set an environment variable in the container. This uses the same format as
[services in systemd](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#Environment=)
and can be listed multiple times.

#### `Exec=`

If this is set then it defines what command line to run in the container. If it is not set the
default entry point of the container image is used. The format is the same as for
[systemd command lines](https://www.freedesktop.org/software/systemd/man/systemd.service.html#Command%20lines).

#### `User=`

The (numeric) uid to run as inside the container. This does not need to match the uid on the host,
which can be set with `HostUser`, but if that is not specified, this uid is also used on the host.

#### `HostUser=`

The host uid (numeric or a username) to run the container as. If this differs from the uid in `User`,
then user namespaces are used to map the ids. If unspecified, this defaults to what was specified in `User`.

#### `Group=`

The (numeric) gid to run as inside the container. This does not need to match the gid on the host,
which can be set with `HostGroup`, but if that is not specified, this gid is also used on the host.

#### `HostGroup=`

The host gid (numeric or group name) to run the container as. If this differs from the gid in `Group`,
then user namespaces are used to map the ids. If unspecified, this defaults to what was specified in `Group`.

#### `NoNewPrivileges=` (defaults to `yes`)

If enabled (which is the default), this disables the container processes from gaining additional privileges via things like
setuid and file capabilities.

#### `DropCapability=` (defaults to `all`)

Drop these capabilities from the default container capability set. The default is `all`, allowing
addition of capabilities with `AddCapability`. Set this to empty to drop no capabilities.
This can be listed multiple times.

#### `AddCapability=`

By default, the container runs with no capabilities (due to DropCapabilities='all'). If any specific
caps are needed, then add them with this key. For example using `AddCapability=CAP_DAC_OVERRIDE`.
This can be listed multiple times.

#### `ReadOnly=` (defaults to `no`)

If enabled, makes image read-only, with /var/tmp, /tmp and /run a tmpfs (unless disabled by `VolatileTmp=no`).

**NOTE:** Podman will automatically copy any content from the image onto the tmpfs

#### `RemapUsers=` (defaults to `no`)

If this is enabled, then host user and group ids are remapped in the container, such that all the uids
starting at `RemapUidStart` (and gids starting at `RemapGidStart`) in the container are chosen from the
available host uids specified by `RemapUidRanges` (and `RemapGidRanges`).

#### `RemapUidStart=` (defaults to `1`)

If `RemapUsers` is enabled, this is the first uid that is remapped, and all lower uids are mapped
to the equivalent host uid. This defaults to 1 so that the host root uid is in the container, because
this means a lot less file ownership remapping in the container image.

#### `RemapGidStart=` (defaults to `1`)

If `RemapUsers` is enabled, this is the first gid that is remapped, and all lower gids are mapped
to the equivalent host gid. This defaults to 1 so that the host root gid is in the container, because
this means a lot less file ownership remapping in the container image.

#### `RemapUidRanges=`

This specifies a comma-separated list of ranges (like `10000-20000,40000-50000`) of available host
uids to use to remap container uids in `RemapUsers`. Alternatively, it can be a username, which means
the available subuids of that user will be used.

If not specified, the default ranges are chosen as the subuids of the `quadlet` user.

#### `RemapGidRanges=`

This specifies a comma-separated list of ranges (like `10000-20000,40000-50000`) of available host
gids to use to remap container gids in `RemapUsers`. Alternatively, it can be a username, which means
the available subgids of that user will be used.

If not specified, the default ranges are chosen as the subgids of the `quadlet` user.

#### `KeepId=` (defaults to `no`, only works for user units)

If this is enabled, then the user uid will be mapped to itself in the container, otherwise it is
mapped to root. This is ignored for system units.

#### `Notify=` (defaults to `no`)

By default, Podman is run in such a way that the systemd startup notify command is handled by
the container runtime. In other words, the service is deemed started when the container runtime
starts the child in the container. However, if the container application supports
[sd_notify](https://www.freedesktop.org/software/systemd/man/sd_notify.html), then setting
`Notify`to true will pass the notification details to the container allowing it to notify
of startup on its own.

#### `Timezone=` (if unset uses system-configured default)

The timezone to run the container in.

#### `RunInit=` (default to `yes`)

If enabled (and it is by default), the container will have a minimal init process inside the
container that forwards signals and reaps processes.

#### `VolatileTmp=` (default to `yes`)

If enabled (and it is by default), the container will have a fresh tmpfs mounted on `/tmp`.

**NOTE:** Podman will automatically copy any content from the image onto the tmpfs

#### `Volume=`

Mount a volume in the container. This is equivalent to the Podman `--volume` option, and
generally has the form `[[SOURCE-VOLUME|HOST-DIR:]CONTAINER-DIR[:OPTIONS]]`.

As a special case, if `SOURCE-VOLUME` ends with `.volume`, a Podman named volume called
`systemd-$name` will be used as the source, and the generated systemd service will contain
a dependency on the `$name-volume.service`. Such a volume can be automatically be lazily
created by using a `$name.volume` quadlet file.

This key can be listed multiple times.

#### `ExposeHostPort=`

Exposes a port, or a range of ports (e.g. `50-59`), from the host to the container. Equivalent
to the Podman `--expose` option.

This key can be listed multiple times.

#### `PublishPort=`

Exposes a port, or a range of ports (e.g. `50-59`), from the container to the host. Equivalent
to the Podman `--publish` option. The format is similar to the Podman options, which is of
the form `ip:hostPort:containerPort`, `ip::containerPort`, `hostPort:containerPort` or
`containerPort`, where the number of host and container ports must be the same (in the case
of a range).

If the IP is set to 0.0.0.0 or not set at all, the port will be bound on all IPv4 addresses on
the host; use [::] for IPv6.

Note that not listing a host port means that Podman will automatically select one, and it
may be different for each invocation of service. This makes that a less useful option. The
allocated port can be found with the `podman port` command.

This key can be listed multiple times.

#### `PodmanArgs=`

This key contains a list of arguments passed directly to the end of the `podman run` command
in the generated file (right before the image name in the command line). It can be used to
access Podman features otherwise unsupported by the generator. Since the generator is unaware
of what unexpected interactions can be caused by these arguments, is not recommended to use
this option.

The format of this is a space separated list of arguments, which can optionally be individually
escaped to allow inclusion of whitespace and other control characters. This key can be listed
multiple times.

#### `Label=`

Set one or more OCI labels on the container. The format is a list of `key=value` items,
similar to `Environment`.

This key can be listed multiple times.

#### `Annotation=`

Set one or more OCI annotations on the container. The format is a list of `key=value` items,
similar to `Environment`.

This key can be listed multiple times.

### Volume units

Volume files are named with a `.volume` extension and contain a section `[Volume]` describing the
named Podman volume. The generated service is a one-time command that ensures that the volume
exists on the host, creating it if needed.

For a volume file named `$NAME.volume`, the generated Podman volume will be called `systemd-$NAME`,
and the generated service file `$NAME-volume.service`.

Using volume units allows containers to depend on volumes being automatically pre-created. This is
particularly interesting when using special options to control volume creation, as Podman will
otherwise create volumes with the default options.

Supported keys in `Volume` section are:

#### `User=`

The host (numeric) uid, or user name to use as the owner for the volume

#### `Group=`

The host (numeric) gid, or group name to use as the group for the volume

#### `Label=`

Set one or more OCI labels on the volume. The format is a list of
`key=value` items, similar to `Environment`.

This key can be listed multiple times.


Set one or more OCI labels on the volume. The format is a list of `key=value` items,
similar to `Environment`.

This key can be listed multiple times.

## EXAMPLES

Example `test.container`:

```
[Unit]
Description=A minimal container
Before=local-fs.target

[Container]
# Use the centos image
Image=quay.io/centos/centos:latest
Volume=test.volume:/data

# In the container we just run sleep
Exec=sleep 60

[Service]
# Restart service when sleep finishes
Restart=always

[Install]
# Start by default on boot
WantedBy=multi-user.target default.target
```

Example `test.volume`:

```
[Volume]
User=root
Group=projectname
Label=org.test.Key=value
```

## SEE ALSO
**[systemd.unit(5)](https://www.freedesktop.org/software/systemd/man/systemd.unit.html)**,
**[systemd.service(5)](https://www.freedesktop.org/software/systemd/man/systemd.service.html)**,
**[podman-run(1)](podman-run.1.md)**
