% podman-info(1)

## NAME
podman\-info - Displays Podman related system information

## SYNOPSIS
**podman info** [*options*]

**podman system info** [*options*]

## DESCRIPTION

Displays information pertinent to the host, current storage stats, configured container registries, and build of podman.


## OPTIONS

**-D**, **--debug**

Show additional information

**-f**, **--format**=*format*

Change output format to "json" or a Go template.


## EXAMPLE

Run podman info with plain text response:
```
$ podman info
host:
  BuildahVersion: 1.4-dev
  Conmon:
    package: Unknown
    path: /usr/libexec/podman/conmon
    version: 'conmon version 1.12.0-dev, commit: d724f3d54ad2d95b6de741085d4990190ebfd7ff'
  Distribution:
    distribution: fedora
    version: "28"
  MemFree: 1271083008
  MemTotal: 33074233344
  OCIRuntime:
    package: runc-1.0.0-51.dev.gitfdd8055.fc28.x86_64
    path: /usr/bin/runc
    version: 'runc version spec: 1.0.0'
  SwapFree: 34309664768
  SwapTotal: 34359734272
  arch: amd64
  cpus: 8
  hostname: localhost.localdomain
  kernel: 4.18.7-200.fc28.x86_64
  os: linux
  uptime: 218h 49m 33.66s (Approximately 9.08 days)
registries:
  docker.io:
    Blocked: true
    Insecure: true
    Location: docker.io
    MirrorByDigestOnly: false
    Mirrors:
    - Insecure: true
      Location: example2.io/example/ubi8-minimal
    Prefix: docker.io
  redhat.com:
    Blocked: false
    Insecure: false
    Location: registry.access.redhat.com/ubi8
    MirrorByDigestOnly: true
    Mirrors:
    - Insecure: false
      Location: example.io/example/ubi8-minimal
    - Insecure: true
      Location: example3.io/example/ubi8-minimal
    Prefix: redhat.com
store:
  ConfigFile: /etc/containers/storage.conf
  ContainerStore:
    number: 37
  GraphDriverName: overlay
  GraphOptions:
  - overlay.mountopt=nodev
  - overlay.override_kernel_check=true
  GraphRoot: /var/lib/containers/storage
  GraphStatus:
    Backing Filesystem: extfs
    Native Overlay Diff: "true"
    Supports d_type: "true"
  ImageStore:
    number: 17
  RunRoot: /var/run/containers/storage

```
Run podman info with JSON formatted response:
```
{
    "host": {
        "BuildahVersion": "1.4-dev",
        "Conmon": {
            "package": "Unknown",
            "path": "/usr/libexec/podman/conmon",
            "version": "conmon version 1.12.0-dev, commit: d724f3d54ad2d95b6de741085d4990190ebfd7ff"
        },
        "Distribution": {
            "distribution": "fedora",
            "version": "28"
        },
        "MemFree": 1204109312,
        "MemTotal": 33074233344,
        "OCIRuntime": {
            "package": "runc-1.0.0-51.dev.gitfdd8055.fc28.x86_64",
            "path": "/usr/bin/runc",
            "version": "runc version spec: 1.0.0"
        },
        "SwapFree": 34309664768,
        "SwapTotal": 34359734272,
        "arch": "amd64",
        "cpus": 8,
        "hostname": "localhost.localdomain",
        "kernel": "4.18.7-200.fc28.x86_64",
        "os": "linux",
        "uptime": "218h 50m 35.02s (Approximately 9.08 days)"
    },
    "insecure registries": {
        "registries": []
    },
    "registries": {
        "registries": [
            "quay.io",
            "registry.fedoraproject.org",
            "docker.io",
            "registry.access.redhat.com"
        ]
    },
    "store": {
        "ContainerStore": {
            "number": 37
        },
        "GraphDriverName": "overlay",
        "GraphOptions": [
            "overlay.mountopt=nodev",
            "overlay.override_kernel_check=true"
        ],
        "GraphRoot": "/var/lib/containers/storage",
        "GraphStatus": {
            "Backing Filesystem": "extfs",
            "Native Overlay Diff": "true",
            "Supports d_type": "true"
        },
        "ImageStore": {
            "number": 17
        },
        "RunRoot": "/var/run/containers/storage"
    }
}
```
Run podman info and only get the registries information.
```
$ podman info --format={{".Registries"}}
map[registries:[docker.io quay.io registry.fedoraproject.org registry.access.redhat.com]]
```

## SEE ALSO
podman(1), containers-registries.conf(5), containers-storage.conf(5)
