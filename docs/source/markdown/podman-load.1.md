% podman-load(1)

## NAME
podman\-load - Load an image from a container image archive into container storage

## SYNOPSIS
**podman load** [*options*] [*name*[:*tag*]]

**podman image load** [*options*] [*name*[:*tag*]]

## DESCRIPTION
**podman load** loads an image from either an **oci-archive** or **docker-archive** stored on the local machine into container storage. **podman load** reads from stdin by default or a file if the **input** option is set.
You can also specify a name for the image if the archive does not contain a named reference, of if you want an additional name for the local image.

The **quiet** option suppresses the progress output when set.
Note: `:` is a restricted character and cannot be part of the file name.


**podman [GLOBAL OPTIONS]**

**podman load [GLOBAL OPTIONS]**

**podman load [OPTIONS] NAME[:TAG]**

## OPTIONS

**--input**, **-i**=*input*

Read from archive file, default is STDIN.

The remote client requires the use of this option.

NOTE: Use the environment variable `TMPDIR` to change the temporary storage location of container images. Podman defaults to use `/var/tmp`.

**--quiet**, **-q**

Suppress the progress output

**--help**, **-h**

Print usage statement

## EXAMPLES

```
$ podman load --quiet -i fedora.tar
```

```
$ podman load -q -i fedora.tar
```

```
$ podman load < fedora.tar
Getting image source signatures
Copying blob sha256:5bef08742407efd622d243692b79ba0055383bbce12900324f75e56f589aedb0
 0 B / 4.03 MB [---------------------------------------------------------------]
Copying config sha256:7328f6f8b41890597575cbaadc884e7386ae0acc53b747401ebce5cf0d624560
 0 B / 1.48 KB [---------------------------------------------------------------]
Writing manifest to image destination
Storing signatures
Loaded image:  registry.fedoraproject.org/fedora:latest
```

```
$ cat fedora.tar | podman load
Getting image source signatures
Copying blob sha256:5bef08742407efd622d243692b79ba0055383bbce12900324f75e56f589aedb0
 0 B / 4.03 MB [---------------------------------------------------------------]
Copying config sha256:7328f6f8b41890597575cbaadc884e7386ae0acc53b747401ebce5cf0d624560
 0 B / 1.48 KB [---------------------------------------------------------------]
Writing manifest to image destination
Storing signatures
Loaded image:  registry.fedoraproject.org/fedora:latest
```

## SEE ALSO
podman(1), podman-save(1), podman-tag(1)

## HISTORY
July 2017, Originally compiled by Urvashi Mohnani <umohnani@redhat.com>
