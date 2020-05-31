% podman-mount(1)

## NAME
podman\-mount - Mount a working container's root filesystem

## SYNOPSIS
**podman mount** [*container* ...]

**podman container mount** [*container* ...]

## DESCRIPTION
Mounts the specified containers' root file system in a location which can be
accessed from the host, and returns its location.

If you execute the command without any arguments, the tool will list all of the
currently mounted containers.

## RETURN VALUE
The location of the mounted file system.  On error an empty string and errno is
returned.

## OPTIONS

**--all**, **-a**

Mount all containers.

**--format**=*format*

Print the mounted containers in specified format (json)

**--latest**, **-l**

Instead of providing the container name or ID, use the last created container.
If you use methods other than Podman to run containers such as CRI-O, the last
started container could be from either of those methods.

The latest option is not supported on the remote client.

**--notruncate**

Do not truncate IDs in output.

## EXAMPLE

```
podman mount c831414b10a3

/var/lib/containers/storage/overlay/f3ac502d97b5681989dff84dfedc8354239bcecbdc2692f9a639f4e080a02364/merged
```

```
podman mount

c831414b10a3 /var/lib/containers/storage/overlay/f3ac502d97b5681989dff84dfedc8354239bcecbdc2692f9a639f4e080a02364/merged
a7060253093b /var/lib/containers/storage/overlay/0ff7d7ca68bed1ace424f9df154d2dd7b5a125c19d887f17653cbcd5b6e30ba1/merged
```
```
podman mount c831414b10a3 a7060253093b

/var/lib/containers/storage/overlay/f3ac502d97b5681989dff84dfedc8354239bcecbdc2692f9a639f4e080a02364/merged
/var/lib/containers/storage/overlay/0ff7d7ca68bed1ace424f9df154d2dd7b5a125c19d887f17653cbcd5b6e30ba1/merged
```

```
podman mount

c831414b10a3 /var/lib/containers/storage/overlay/f3ac502d97b5681989dff84dfedc8354239bcecbdc2692f9a639f4e080a02364/merged
a7060253093b /var/lib/containers/storage/overlay/0ff7d7ca68bed1ace424f9df154d2dd7b5a125c19d887f17653cbcd5b6e30ba1/merged
```

## SEE ALSO
podman(1), podman-umount(1), mount(8)
