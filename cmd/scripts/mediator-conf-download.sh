#!/bin/bash
# A helper to download mediator-client.json from the SC pod.
# Since the "tos" command needs high privileges, this script expects to run as root.

set -e
test -n "${1}"
tmp_dir=$(mktemp -t -d "mediator.XXXXXXXX")
trap "rm -rf '${tmp_dir}'" EXIT
touch "${tmp_dir}/mediator-client.json"
/usr/local/bin/tos scripts sc pull "${tmp_dir}" --overwrite mediator-client.json || exit 0
chown tufin-admin:tufin-admin "${tmp_dir}/mediator-client.json"
mv "${tmp_dir}/mediator-client.json" "${1}"
