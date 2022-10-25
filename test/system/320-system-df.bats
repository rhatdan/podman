#!/usr/bin/env bats   -*- bats -*-
#
# tests for podman system df
#

load helpers

function teardown() {
    basic_teardown

    # In case the active-volumes test failed: clean up stray volumes
    run_podman volume rm -a
}

@test "podman system df - basic functionality" {
    run_podman system df
    is "$output" ".*Images  *1 *0 "       "Exactly one image"
    is "$output" ".*Containers *0 *0 "    "No containers"
    is "$output" ".*Local Volumes *0 *0 " "No volumes"
}

@test "podman system df - with active containers and volumes" {
    run_podman run    -v /myvol1 --name c1 $IMAGE true
    run_podman run -d -v /myvol2 --name c2 $IMAGE \
               sh -c 'while ! test -e /stop; do sleep 0.1;done'

    run_podman system df --format '{{ .Type }}:{{ .Total }}:{{ .Active }}'
    is "${lines[0]}" "Images:1:1"        "system df : Images line"
    is "${lines[1]}" "Containers:2:1"    "system df : Containers line"
    is "${lines[2]}" "Local Volumes:2:2" "system df : Volumes line"

    # Try -v. (Grrr. No way to specify individual formats)
    #
    # Yes, I know this would be more elegant as a separate @test, but
    # container/volume setup/teardown costs ~3 seconds and that matters.
    run_podman system df -v
    is "${lines[2]}" \
       "${PODMAN_TEST_IMAGE_REGISTRY}/${PODMAN_TEST_IMAGE_USER}/${PODMAN_TEST_IMAGE_NAME} * ${PODMAN_TEST_IMAGE_TAG} [0-9a-f]* .* 2" \
       "system df -v: the 'Images' line"

    # Containers are listed in random order. Just check that each has 1 volume
    is "${lines[5]}" \
       "[0-9a-f]\{12\} *[0-9a-f]\{12\} .* 1 .* c[12]" \
       "system df -v, 'Containers', first line"
    is "${lines[6]}" \
       "[0-9a-f]\{12\} *[0-9a-f]\{12\} .* 1 .* c[12]" \
       "system df -v, 'Containers', second line"

    # Volumes, likewise: random order.
    is "${lines[9]}" "[0-9a-f]\{64\} *[01] * 0B" \
       "system df -v, 'Volumes', first line"
    is "${lines[10]}" "[0-9a-f]\{64\} *[01] * 0B" \
       "system df -v, 'Volumes', second line"

    # Make sure that the return image "raw" size is correct
    run_podman image inspect $IMAGE --format "{{.Size}}"
    expectedSize="$output"

    run_podman system df --format "{{.RawSize}}"
    is "${lines[0]}" "$expectedSize" "raw image size is correct"

    # Clean up and check reclaimable image data
    run_podman system df --format '{{.Reclaimable}}'
    is "${lines[0]}" "0B (0%)" "cannot reclaim image data as it's still used by the containers"

    run_podman exec c2 touch /stop
    run_podman wait c2

    # Create a second image by committing a container.
    run_podman container commit -q c1
    image="$output"

    run_podman system df --format '{{.Reclaimable}}'
    is "${lines[0]}" ".* (100%)" "100 percent of image data is reclaimable because $IMAGE has unique size of 0"

    # Make sure the unique size is now really 0.  We cannot use --format for
    # that unfortunately but we can exploit the fact that $IMAGE is used by
    # two containers.
    run_podman system df -v
    is "$output" ".*0B\\s\\+2.*"

    run_podman rm c1 c2

    run_podman system df --format '{{.Reclaimable}}'
    is "${lines[0]}" ".* (100%)" "100 percent of image data is reclaimable because all containers are gone"

    run_podman rmi $image
    run_podman volume rm -a
}

# vim: filetype=sh
