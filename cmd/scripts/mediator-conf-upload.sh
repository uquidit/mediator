#!/bin/bash
# A helper to upload mediator-client.json on the SC pod.
# Since the "tos" command needs high privileges, this script expects to run as root.

set -e
test -n "${1}"
tmp_dir=$(mktemp -t -d "mediator.XXXXXXXX")
trap "rm -rf '${tmp_dir}'" EXIT
cp "${1}" "${tmp_dir}/mediator-client.json"
/usr/local/bin/tos scripts sc push --overwrite "${tmp_dir}/mediator-client.json"
