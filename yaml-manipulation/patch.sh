#!/bin/bash
INPUT="input.yml"

# Store modified YAML to a variable
MODIFIED_YAML_RQ=$(yq eval '
.container_config |= map(
    .resource_quota = {
        "requests.cpu": "1",
        "requests.memory": "1Gi",
        "limits.cpu": "2",
        "limits.memory": "2Gi",
        "pods": "2",
        "persistentvolumeclaims": "5",
        "requests.storage": "5Gi"
    }
)
' "$INPUT")

MODIFIED_YAML=$(yq eval '
.container_config[] |= (
  .namespaces = ["namespace1", "namespace2"]
)
' <<< "$MODIFIED_YAML_RQ")

# Print or use
echo "$MODIFIED_YAML" > output.yml