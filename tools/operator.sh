#!/bin/bash
set -e

current_dir=$1
if [ -z "$1" ]; then current_dir=.; fi

operator_cert=$(<"$current_dir"/operator.crt)
operator_key=$(<"$current_dir"/operator.key)
root_ca=$(<"$current_dir"/rootCA.crt)

cat > ../operator/certs.go <<EOF
package main

const clientCert = \`
${operator_cert}\`

const clientKey = \`
${operator_key}\`

const rootCA = \`
${root_ca}\`
EOF
