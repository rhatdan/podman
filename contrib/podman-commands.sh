#!/bin/bash

# override with, e.g.,  PODMAN=./bin/podman
PODMAN=${PODMAN:-podman}

function die() {
    echo "FATAL: $*" >&2
    exit 1
}


function podman_commands() {
    $PODMAN help "$@" |\
        awk '/^Available Commands:/{ok=1;next}/^Flags:/{ok=0}ok { print $1 }' |\
        grep .
}

function podman_man() {
    # - tr: remove nonalphanumeric and nonwhitespace characters
    # - awk: print everything from COMMANDS or SUBCOMMANDS to next section
    # - egrep: print the first or second column, matching podman-xxxxx
    # - sed: strip away all except the last word, and remove the '(1)'
    man ./docs/$1.1 |\
        tr -cd 'a-zA-Z0-9 \012-' |\
        awk '/^(SUB)?COMMANDS/{ok=1;next}/^[^ ]/{ok=0}ok' |\
        egrep -o 'podman-[a-z0-9-]+' |\
        sed -e 's/^.*-//' -e 's/1$//'
}

function compare_help_and_man() {
    echo
    echo "checking: $@"

    # e.g. podman, podman-image, podman-volume
    basename=$(echo podman "$@" | sed -e 's/ /-/g')

    # uniq is necessary because, e.g., in podman-image.1 the 'list' and 'ls'
    # commands reference 'podman-images(1)' (?!?)
    podman_commands "$@" | sort | uniq > /tmp/${basename}_help.txt
    podman_man $basename | sort | uniq > /tmp/${basename}_man.txt

    diff -u /tmp/${basename}_help.txt /tmp/${basename}_man.txt

    # Now look for subcommands, e.g. container, image
    for cmd in $(< /tmp/${basename}_help.txt); do
        # FIXME FIXME FIXME
        usage=$($PODMAN "$@" $cmd --help | grep -A2 '^Usage:' | grep . | tail -1)

        # if ends in '[command]', recurse into subcommands
        if expr "$usage" : '.*\[command\]$' >/dev/null; then
            compare_help_and_man "$@" $cmd
        fi
    done
}


compare_help_and_man
