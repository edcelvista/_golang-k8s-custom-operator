#!/bin/bash
kubectl get nodes -o custom-columns=NAME:.metadata.name --no-headers | while read node; do
  echo "Node: $node"

  allocatable_cpu=$(kubectl get node "$node" -o jsonpath='{.status.allocatable.cpu}')
  allocatable_mem=$(kubectl get node "$node" -o json | jq -r '.status.allocatable.memory | capture("(?<value>[0-9]+)(?<unit>Ki|Mi|Gi)?")
    | (.value | tonumber) *
      (if .unit == "Gi" then 1
       elif .unit == "Mi" then 1 / 1024
       elif .unit == "Ki" then 1 / (1024 * 1024)
       else 1 / (1024 * 1024)
       end)')

  requested_cpu=$(kubectl get pods --all-namespaces --field-selector spec.nodeName=$node \
    -o json | jq ' [
    .items[].spec.containers[].resources.requests.cpu // "0"
    | capture("(?<value>[0-9]+)(?<unit>m)?")
    | (.value | tonumber) *
      (if .unit == "m" then 1 / 1000
       else 1
       end)
  ] | add ')

  requested_mem=$(kubectl get pods --all-namespaces --field-selector spec.nodeName=$node \
    -o json | jq ' [
    .items[].spec.containers[].resources.requests.memory // "0"
    | capture("(?<value>[0-9]+)(?<unit>Ki|Mi|Gi)?")
    | (.value | tonumber) *
      (if .unit == "Gi" then 1
       elif .unit == "Mi" then 1 / 1024
       elif .unit == "Ki" then 1 / (1024 * 1024)
       else 1 / (1024 * 1024)
       end)
  ] | add')

  echo "  Allocatable CPU: $allocatable_cpu"
  echo "  Requested CPU: $requested_cpu"
  echo "  Allocatable Memory: $allocatable_mem"
  echo "  Requested Memory: ${requested_mem}"
  echo
done