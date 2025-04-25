#!/bin/bash
kubectl create secret tls certs --key tls.key --cert tls.crt --dry-run=client -o yaml