#!/bin/bash
set -e

current_dir=$1
if [ -z "$1" ]; then current_dir=.; fi

server_cert=$(<"$current_dir"/server.crt)
server_key=$(<"$current_dir"/server.key)
root_ca=$(<"$current_dir"/rootCA.crt)

cat > ../server/certs.go <<EOF
package main

const serverCert = \`
${server_cert}\`

const serverKey = \`
${server_key}\`

const rootCA = \`
${root_ca}\`
EOF
