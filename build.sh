#!/bin/bash

set -e

function log(){
    local lvl="${1}" ; shift ;
    local msg="${*}"
    local now=$(date +"%Y-%m-%d %H:%M:%S")
    echo "[${now}][${lvl}] ${msg}"
}

function cleanup(){
    if [ ! -z "${tmp_dir}" ] && [ -e "${tmp_dir}" ]; then
        log "DEBUG" "Removing temp dir: ${tmp_dir}"
        rm -rf "${tmp_dir}"
    fi
}

function fatal(){
    log "FATAL" "${*}"
    exit 1
}

trap 'fatal "Caught Ctrl-C"' INT
trap 'cleanup' EXIT

# Make the version just the date
# It's an useful info, yet easy to retrieve
version=$(date '+%Y%m%d%H%M%S')

# Generate some random strings
salt=$(tr -cd '[:alnum:]' < /dev/urandom | head -c16)
pepper=$(tr -cd '[:alnum:]' < /dev/urandom | head -c16)
key=$(tr -cd '[:alnum:]' < /dev/urandom | head -c16)

# TOTP secrets need to be base32 encoded
s1=$(tr -cd '[:alnum:]' < /dev/urandom | head -c15 | base32)
s2=$(tr -cd '[:alnum:]' < /dev/urandom | head -c15 | base32)

# Get a directory to work in
tmp_dir=$(mktemp -t -d "mediator.XXXXXXXX")
log "INFO" "Working in ${tmp_dir}"

out="${tmp_dir}/mediator"
mkdir "${out}"

# Ask Go for staticaly-linked binaries
export CGO_ENABLED=0

# Build the binaries
cd cmd
for elt in mediator-client mediator-server mediator-cli; do
    log "INFO" "Building ${elt}"

    cd "${elt}"
    go build -o "${out}" --ldflags="\
        -X 'main.Version=${version}'\
        -X 'uqtu/mediator/mediatorscript.salt=${salt}' \
        -X 'uqtu/mediator/mediatorscript.pepper=${pepper}' \
        -X 'uqtu/mediator/mediatorscript.secretKey=${key}' \
        -X 'uqtu/mediator/totp.secretMS1=${s1}' \
        -X 'uqtu/mediator/totp.secretMS2=${s2}' \
        " .
    cd ..
done
cd ..

log "INFO" "Copying extra files"

# Copy the mediator-(client|server)_dist.yml files
cp ./cmd/mediator-server/mediator-server_dist.yml "${out}"
cp ./cmd/mediator-client/mediator-client_dist.yml "${out}"

# Generate another random string, use it as the webserver secret
s3=$(tr -cd '[:alnum:]' < /dev/urandom | head -c32)
sed -i "s/EnterYourSecretHere/${s3}/" "${out}/mediator-server_dist.yml"

# Copy the systemctl service file
cp ./cmd/mediator-server/mediator-server.service "${out}"

# Copy the scripts
cp -r ./cmd/scripts/ "${out}/"

log "INFO" "Building the final archive"

# Make a nice archive
tar -czf "${tmp_dir}/mediator.tar.gz" -C "${tmp_dir}" "mediator"

log "INFO" "Moving the archive to ${PWD}"

# Move the archive to ./
mv "${tmp_dir}/mediator.tar.gz" ./
