#!/bin/sh
podman --help | sed -n -Ee '0,/Available Commands/d' -e '/^Flags/q;p' | sed '/^$/d' | awk '{ print $1 }' > /tmp/podman.cmd
man podman | sed -n -e '0,/COMMANDS/d' -e '/^FILES/q;p' | grep podman | cut -f2 -d- | cut -f1 -d\( > /tmp/podman.man
echo diff -B -u /tmp/podman.cmd /tmp/podman.man
diff -B -u /tmp/podman.cmd /tmp/podman.man

podman image --help | sed -n -e '0,/Available Commands/d' -e '/^Flags/q;p' | sed '/^$/d' | awk '{ print $1 }' > /tmp/podman-image.cmd
man podman image | sed -n -Ee '0,/COMMANDS/d'  -e 's/^[[:space:]]*//' -e '/^SEE ALSO/q;p' | grep podman | cut -f1 -d' ' | sed 's/^.//' > /tmp/podman-image.man
echo diff -B -u /tmp/podman-image.cmd /tmp/podman-image.man
diff -B -u /tmp/podman-image.cmd /tmp/podman-image.man

podman container --help | sed -n -e '0,/Available Commands/d' -e '/^Flags/q;p' | sed '/^$/d' | awk '{ print $1 }' > /tmp/podman-container.cmd
man podman container | sed -n -Ee '0,/COMMANDS/d'  -e 's/^[[:space:]]*//' -e '/^SEE ALSO/q;p' | grep podman | cut -f1 -d' ' | sed 's/^.//' > /tmp/podman-container.man
echo diff -B -u /tmp/podman-container.cmd /tmp/podman-container.man
diff -B -u /tmp/podman-container.cmd /tmp/podman-container.man
