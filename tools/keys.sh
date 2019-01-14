#!/bin/bash
set -e

# Generate root CA
openssl genrsa -out rootCA.key 4096
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -subj '/CN=porth Root CA/O=hodo.ro/C=RO' -out rootCA.crt

# Generate server certificate

openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj '/CN=porth Server/O=hodo.ro/C=RO' -out server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost,DNS:porth.hodo.ro") -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out server.crt -days 365 -sha256

# Generate client certificate

openssl genrsa -out client.key 2048
openssl req -new -key client.key -subj '/CN=porth Client/O=hodo.ro/C=RO' -out client.csr
openssl x509 -req -in client.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out client.crt -days 365 -sha256

# Generate operator certificate

openssl genrsa -out operator.key 2048
openssl req -new -key operator.key -subj '/CN=porth Operator/O=hodo.ro/C=RO' -out operator.csr
openssl x509 -req -in operator.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out operator.crt -days 365 -sha256

sleep 1s

# Generate certs.go files
./server.sh
./client.sh
./operator.sh

