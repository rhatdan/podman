#!/usr/bin/env bats

load helpers

@test "podman info - basic test" {
    skip_if_remote "capitalization inconsistencies"

    run_podman info

    expected_keys="
buildahversion: *[0-9.]\\\+
conmon:\\\s\\\+package:
distribution:
ociruntime:\\\s\\\+name:
os:
rootless:
registries:
store:
graphdrivername:
graphroot:
graphstatus:
imagestore:\\\s\\\+number: 1
runroot:
"
    while read expect; do
        is "$output" ".*$expect" "output includes '$expect'"
    done < <(parse_table "$expected_keys")
}

@test "podman info - json" {
    skip_if_remote "capitalization inconsistencies"

    run_podman info --format=json

    expr_nvr="[a-z0-9-]\\\+-[a-z0-9.]\\\+-[a-z0-9]\\\+\."
    expr_path="/[a-z0-9\\\/.-]\\\+\\\$"

    tests="
host.buildahVersion       | [0-9.]
host.conmon.path          | $expr_path
host.ociRuntime.path      | $expr_path
store.configFile          | $expr_path
store.graphDriverName     | [a-z0-9]\\\+\\\$
store.graphRoot           | $expr_path
store.imageStore.number   | 1
"

    parse_table "$tests" | while read field expect; do
        actual=$(echo "$output" | jq -r ".$field")
        dprint "# actual=<$actual> expect=<$expect>"
        is "$actual" "$expect" "jq .$field"
    done

}

# vim: filetype=sh
