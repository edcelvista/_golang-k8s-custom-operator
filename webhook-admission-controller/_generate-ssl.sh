#!/bin/bash
# openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=*.custom-webhook.svc/O=edcelvistadotcom"
openssl req -new -newkey rsa:2048 -days 365 -nodes -keyout tls.key -out tls.csr -config openssl.cnf
openssl x509 -req -in tls.csr -signkey tls.key -out tls.crt -extensions v3_req -extfile openssl.cnf