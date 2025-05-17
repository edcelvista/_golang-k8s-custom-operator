# CREATE CUSTOM OPERATOR & CRD
![alt text](image.png)

## Deployment Checker Operator (Namespaced Scoped)
Check and Reconcile app deployments. Create if not exists or Update image tag if outdated based on CRD Resource.

### CRD
```
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: myapps.k8s.edcelvista.com
spec:
  group: k8s.edcelvista.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                image:
                  type: string
                replicas:
                  type: integer
                appSelector: 
                  type: string
  scope: Namespaced
  names:
    plural: myapps
    singular: myapp
    kind: MyApp
    shortNames:
    - ma
```
### CRD Resource
```
apiVersion: "k8s.edcelvista.com/v1"
kind: MyApp
metadata:
  name: myapp-resource
  namespace: demo
spec:
  image: edcelvista/ubuntu24-network-tools:53
  replicas: 1
  appSelector: myapp-resource
```
**Note:** Controller scan, check and reconcile app deployments in demo and identify the deployments via label `appSelector: myapp-resource`
[Complete Flow Diagram](https://raw.githubusercontent.com/edcelvista/_golang-k8s-custom-operator/refs/heads/main/flow-diagram.draw.io.drawio)
### Role
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: custom-controller-role
rules:
- apiGroups: ["k8s.edcelvista.com"]
  resources: ["MyApp"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["deployment"]
  verbs: ["create", "get", "update", "list", "watch", "delete"]
```
### Operator
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: custom-operator
  labels:
    app: custom-operator
spec:
  selector:
    matchLabels:
      app: custom-operator
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: custom-operator
    spec:
      containers:
        - name: custom-operator
          image: edcelvista/ubuntu24-network-tools:53
          imagePullPolicy: Always
          command: ["custom-operator-deployment-recon-linux"]
          resources: 
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 100Mi
          envFrom:
            - configMapRef:
                name: operator-cm
          volumeMounts:
          - name: config-volume
            mountPath: "/opt/config"
            readOnly: true
          - name: localtime
            mountPath: /etc/localtime
      volumes:
        - name: config-volume
          configMap:
            name: kubeconfig
        - name: localtime
          hostPath:
            path: /usr/share/zoneinfo/Asia/Shanghai
      restartPolicy: Always
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kubeconfig
data:
  config: |
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: LS0tLS...
        server: https://172.31.24.226:6443
      name: kubernetes-aws
    contexts:
    - context:
        cluster: kubernetes-aws
        user: kubernetes-aws-admin
      name: kubernetes-aws-admin@kubernetes
    current-context: kubernetes-aws-admin@kubernetes
    kind: Config
    preferences: {}
    users:
    - name: kubernetes-aws-admin
      user:
        client-certificate-data: LS0tLS...
        client-key-data: LS0tLS...
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: operator-cm
data:
  CUSTOM_KUBE_CONFIG_PATH: "/opt/config/config"
  CRDNAME: "myapps.k8s.edcelvista.com"
  CRDGROUP: "k8s.edcelvista.com"
  CRDRESOURCE: "myapps"
  CRDVERSION: "v1"
  TARGETNAMESPACE: "demo"
  APPNAME: "demo-crd-resource"
  APPJSONTEMPLATE: "{ \"apiVersion\": \"apps/v1\", \"kind\": \"Deployment\", \"spec\": { \"selector\": { \"matchLabels\": {} }, \"template\": { \"spec\": { \"containers\": [ { \"command\": [ \"sleep\", \"infinity\" ], \"resources\": {}, \"terminationMessagePath\": \"/dev/termination-log\", \"terminationMessagePolicy\": \"File\", \"imagePullPolicy\": \"IfNotPresent\" } ], \"restartPolicy\": \"Always\", \"terminationGracePeriodSeconds\": 30, \"dnsPolicy\": \"ClusterFirst\", \"securityContext\": {}, \"schedulerName\": \"default-scheduler\" } }, \"strategy\": { \"type\": \"RollingUpdate\", \"rollingUpdate\": { \"maxUnavailable\": \"25%\", \"maxSurge\": \"25%\" } }, \"revisionHistoryLimit\": 10, \"progressDeadlineSeconds\": 600 } }"
  INVERVAL: "30"
```
**Note:** Operator controls and dictates the base app deployment structure.

---

## Secret Checker Operator (Cluster Scoped)
Check and Reconcile app secrets. Create if not exists or Update secret if outdated based on CRD Resource.

### CRD
```
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: sectools.k8s.edcelvista.com
spec:
  group: k8s.edcelvista.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                secret:
                  type: object
                  properties:
                    type:
                      type: string
                    name:
                      type: string
                    data:
                      type: object
                      description: Base64-encoded data entries like in Kubernetes Secrets.
                      additionalProperties:
                        type: object
                        properties:
                          tls.crt:
                            type: string
                            format: byte  # base64-encoded
                          tls.key:
                            type: string
                            format: byte  # base64-encoded
                    stringData:
                      type: object
                      description: Plaintext string data; will be base64-encoded into `data`
                      additionalProperties:
                        type: string
                  required: 
                    - data
  scope: Cluster
  names:
    plural: sectools
    singular: sectool
    kind: SecTool
    shortNames:
    - st
```
### CRD Resource
```
apiVersion: "k8s.edcelvista.com/v1"
kind: SecTool
metadata:
  name: sectool-resource
spec:
  secret:
    type: kubernetes.io/tls
    name: edcelvistadotcom-aws-tls
    data:
      additionalProperties:
        tls.crt: >-
          LS0tLS...
        tls.key: >-
          LS0tLS...
```
### Cluster Role
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: custom-controller-role-st
rules:
- apiGroups: ["k8s.edcelvista.com"]
  resources: ["SecTool"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["secret"]
  verbs: ["create", "get", "update", "list", "watch", "delete"]
```
### Operator
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: custom-operator
  labels:
    app: custom-operator
spec:
  selector:
    matchLabels:
      app: custom-operator
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: custom-operator
    spec:
      containers:
        - name: custom-operator
          image: edcelvista/ubuntu24-network-tools:53
          imagePullPolicy: Always
          command: ["custom-operator-secret-recon-linux"]
          resources: 
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 100Mi
          envFrom:
            - configMapRef:
                name: operator-cm
          volumeMounts:
          - name: config-volume
            mountPath: "/opt/config"
            readOnly: true
          - name: localtime
            mountPath: /etc/localtime
      volumes:
        - name: config-volume
          configMap:
            name: kubeconfig
        - name: localtime
          hostPath:
            path: /usr/share/zoneinfo/Asia/Shanghai
      restartPolicy: Always
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kubeconfig
data:
  config: |
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: LS0tLS...
        server: https://172.31.24.226:6443
      name: kubernetes-aws
    contexts:
    - context:
        cluster: kubernetes-aws
        user: kubernetes-aws-admin
      name: kubernetes-aws-admin@kubernetes
    current-context: kubernetes-aws-admin@kubernetes
    kind: Config
    preferences: {}
    users:
    - name: kubernetes-aws-admin
      user:
        client-certificate-data: LS0tLS...
        client-key-data: LS0tLS...
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: operator-cm
data:
  CUSTOM_KUBE_CONFIG_PATH: "/opt/config/config"
  CRDNAME: "sectools.k8s.edcelvista.com"
  CRDGROUP: "k8s.edcelvista.com"
  CRDRESOURCE: "sectools"
  CRDVERSION: "v1"
  EXCLUDENAMESPACE: "argocd,castai-agent,custom-operator,default,krakend,kube-bench,kube-flannel,kube-node-lease,kube-public,kube-system,kyverno,monitoring,nginx-ingress"
  APPNAME: "edcelvistadotcom-aws-tls"
  INVERVAL: "30"
```
**Note:** Operator controls and dictates the what namespace to be excluded in secret injection.

---

# Admission Controller - Mutating & Validating
```
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: validating-always-allow
webhooks:
  - name: validating-always-allow.edcelvista.com
    rules:
      - apiGroups:
        - ''
        apiVersions:
        - 'v1'
        operations:
        - CREATE
        resources:
        - 'pods'
        scope: 'Namespaced'
    clientConfig:
      # url: https://webhook.custom-webhook.svc.cluster.local/webhook/validating/always-allow
      service:
        name: webhook
        namespace: custom-webhook
        path: "/webhook/validating/always-allow"
      caBundle: LS0tLS...
    admissionReviewVersions: ["v1"]
    sideEffects: None
    # failurePolicy: Ignore
    matchPolicy: Equivalent
    namespaceSelector: {}
    objectSelector:
      matchExpressions:
      - key: custom-webhook.edcelvista.com/validate-always-allow
        operator: Exists
    timeoutSeconds: 5
---
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: mutating-always-allow
webhooks:
  - name: mutating-always-allow.edcelvista.com
    rules:
      - apiGroups:
        - ''
        apiVersions:
        - 'v1'
        operations:
        - CREATE
        resources:
        - 'pods'
        scope: 'Namespaced'
    clientConfig:
      service:
        name: webhook
        namespace: custom-webhook
        path: /webhook/mutating/always-allow
      caBundle: LS0tLS...
    admissionReviewVersions: ["v1"]
    sideEffects: None
    # failurePolicy: Ignore
    matchPolicy: Equivalent
    namespaceSelector: {}
    objectSelector:
      matchExpressions:
      - key: custom-webhook.edcelvista.com/mutate-always-allow
        operator: Exists
    timeoutSeconds: 5
```
### Webhook API
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webhook-operator
  labels:
    app: webhook-operator
spec:
  selector:
    matchLabels:
      app: webhook-operator
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: webhook-operator
    spec:
      containers:
        - name: webhook-operator
          image: edcelvista/ubuntu24-network-tools:53
          imagePullPolicy: Always
          command: ["webhook-linux"]
          resources: 
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 100Mi
          envFrom:
            - configMapRef:
                name: webhook-cm
          ports:
          - name: webport
            containerPort: 8443
            protocol: TCP
          volumeMounts:
          - name: certs-volume
            mountPath: "/certs"
            readOnly: true
          - name: localtime
            mountPath: /etc/localtime
      volumes:
        - name: certs-volume
          secret:
            secretName: certs
        - name: localtime
          hostPath:
            path: /usr/share/zoneinfo/Asia/Shanghai
      restartPolicy: Always
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: webhook-cm
data:
  PORT: ":8443"
  SSL_CERT: "/certs/tls.crt"
  SSL_KEY: "/certs/tls.key"
  VALIDATE_LABEL: "costCenter,tenantName,supportEmail"
  MUTATE_PATCH: "[{ \"op\": \"add\", \"path\": \"/metadata/labels/mutatedLabel\", \"value\": \"mutatedLabelValue\" }, { \"op\": \"add\", \"path\": \"/spec/containers/0/lifecycle\", \"value\": { \"postStart\": { \"exec\": { \"command\": [ \"/bin/sh\", \"-c\", \"curl -sk --location 'https://webhook.custom-webhook.svc:443/echo' --header 'Content-Type: application/json' --data '{ \\\"postStart\\\": \\\"$POD_NAMESPACE/$POD_NAME with $POD_IP\\\" }'\" ] } }, \"preStop\": { \"exec\": { \"command\": [ \"/bin/sh\", \"-c\", \"curl -sk --location 'https://webhook.custom-webhook.svc:443/echo' --header 'Content-Type: application/json' --data '{ \\\"preStop\\\": \\\"$POD_NAMESPACE/$POD_NAME with $POD_IP\\\" }'\" ] } } } }]"
  IS_DEBUG: "true"
---
apiVersion: v1
kind: Secret
metadata:
  name: certs
type: kubernetes.io/tls
data:
  tls.crt: LS0tLS...
  tls.key: LS0tLS...
---
apiVersion: v1
kind: Service
metadata:
  name: webhook
spec:
  ports:
    - name: webport
      protocol: TCP
      port: 443
      targetPort: 8443
  selector:
    app: webhook-operator
  type: ClusterIP
  sessionAffinity: None
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```
**Note:**  
Webhook handles the validate condition components for Validating ie. required labels and handles Mutating patch definition.
Webhook API [repo here](https://github.com/edcelvista/_golang-k8s-custom-operator/tree/main/webhook-admission-controller/app)