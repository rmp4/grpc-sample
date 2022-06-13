find . -type f -name '*.pem' -delete

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -sha256 -keyout configs/ca-key.pem -out configs/ca-cert.pem -config configs/ssl.conf

echo "CA's self-signed certificate"
openssl x509 -in configs/ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -sha256  -keyout configs/server/server-key.pem -out configs/server/server-req.pem -config configs/server/server.conf

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -sha256 -in configs/server/server-req.pem -days 60 -CA configs/ca-cert.pem -CAkey configs/ca-key.pem -CAcreateserial -out configs/server/server-cert.pem -extfile configs/server/server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in configs/server/server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -sha256 -keyout configs/client/client-key.pem -out configs/client/client-req.pem -config configs/client/client.conf

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -sha256 -in configs/client/client-req.pem -days 60 -CA configs/ca-cert.pem -CAkey configs/ca-key.pem -CAcreateserial -out configs/client/client-cert.pem -extfile configs/client/client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in configs/client/client-cert.pem -noout -text