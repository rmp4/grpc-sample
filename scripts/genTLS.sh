#!/bin/bash

openssl genrsa -out configs/ca.key 4096
# openssl req -new -x509 -days 365 -sha256  -subj "/C=TW/ST=Taipei/O=your_company/OU=your_department/CN=localhost"   -out configs/ca.crt
openssl req -x509 -new -nodes -sha256 -utf8 -days 3650 -key configs/ca.key rsa:2048 -keyout server.key -out configs/ca.crt -config configs/ssl.conf

openssl genrsa -out configs/server/server.key 2048
openssl req -new -key configs/server/server.key -subj "/C=TW/ST=Taipei/O=your_company/OU=your_department/CN=your_service_FQDN" -out configs/server/server.csr

openssl x509 -req -CAcreateserial -days 30 -sha256 -CA configs/ca.crt -CAkey configs/ca.key -in configs/server/server.csr -out configs/server/server.crt



# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout configs/ca.key -out configs/ca.pem -config configs/ssl.conf

echo "CA's self-signed certificate"
openssl x509 -in configs/ca.pem  -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -config configs/ssl.conf

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text


