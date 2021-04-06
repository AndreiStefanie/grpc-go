#!/bin/bash

# Generate the CA certificate and private key
openssl req -x509 \
  -newkey rsa:4096 -days 365 -nodes \
  -keyout ssl/ca-key.pem \
  -out ssl/ca-cert.pem \
  -subj "/C=RO/CN=localhost"

# Generate the server key and CSR
openssl req -newkey rsa:4096 -nodes \
  -keyout ssl/server-key.pem \
  -out ssl/server-req.pem \
  -subj "/C=RO/CN=localhost"

# Sign the server CSR with the CA private key
openssl x509 -req -in ssl/server-req.pem \
  -days 60 \
  -CA ssl/ca-cert.pem -CAkey ssl/ca-key.pem -CAcreateserial \
  -out ssl/server-cert.pem \
  -extfile ssl/server-ext.cnf
