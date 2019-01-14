#!/bin/bash
set -e

current_dir=$1
if [ -z "$1" ]; then current_dir=.; fi

client_cert=$(<"$current_dir"/client.crt)
client_key=$(<"$current_dir"/client.key)
root_ca=$(<"$current_dir"/rootCA.crt)

cat > ../client/certs.go <<EOF
package main

const clientCert = \`
${client_cert}\`

const clientKey = \`
${client_key}\`

const rootCA = \`
${root_ca}\`
EOF
