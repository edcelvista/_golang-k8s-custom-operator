#!/bin/bash
KUBESA=kubeadmin-demo
KUBENS=demo
SECRET_NAME=$(kubectl get sa $KUBESA -n $KUBENS -o jsonpath='{.secrets[0].name}')
TOKEN=$(kubectl get secret $SECRET_NAME -n $KUBENS -o jsonpath='{.data.token}' | base64 -d)

# if token not generated automatically
kubectl create secret generic $KUBESA-token-secret \
  --from-literal=token=$TOKEN \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl patch serviceaccount $KUBESA -n $KUBENS \
  -p '{"secrets":[{"name":"$KUBESA-token-secret"}]}'