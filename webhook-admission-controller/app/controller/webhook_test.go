package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebhookValidatingHandlerPOSTPod(t *testing.T) {
	jsonData := []byte(`{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1","request":{"uid":"cb2373d5-7933-43a4-9ca7-0be7db7c6f33","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"name":"demo-crd-resource-76bf5f8695-sf4bt","namespace":"demo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"50211333-a14e-4c1a-9c38-e00d5b1af517","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"],"extra":{"authentication.kubernetes.io/credential-id":["JTI=7edabfc3-6128-44f4-8d8e-bc35d1ccd47d"]}},"object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"demo-crd-resource-76bf5f8695-sf4bt","generateName":"demo-crd-resource-76bf5f8695-","namespace":"demo","uid":"7df01e81-bb71-49d5-a36b-74de72836d75","creationTimestamp":"2025-04-29T04:39:35Z","labels":{"appSelector":"myapp-resource","costCenter":"test","custom-webhook.edcelvista.com/validate-always-allow":"true","pod-template-hash":"76bf5f8695","supportEmail":"administrator_edcelvista.com","tenantName":"test"},"annotations":{"kubectl.kubernetes.io/restartedAt":"2025-04-29T04:39:32Z"},"ownerReferences":[{"apiVersion":"apps/v1","kind":"ReplicaSet","name":"demo-crd-resource-76bf5f8695","uid":"ee15fb57-668f-4d47-b710-29e8c37fbc7b","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"v1","time":"2025-04-29T04:39:35Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/restartedAt":{}},"f:generateName":{},"f:labels":{".":{},"f:appSelector":{},"f:costCenter":{},"f:custom-webhook.edcelvista.com/validate-always-allow":{},"f:pod-template-hash":{},"f:supportEmail":{},"f:tenantName":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"ee15fb57-668f-4d47-b710-29e8c37fbc7b\"}":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demo-crd-resource\"}":{".":{},"f:command":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:enableServiceLinks":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}]},"spec":{"volumes":[{"name":"kube-api-access-5drvz","projected":{"sources":[{"serviceAccountToken":{"expirationSeconds":3607,"path":"token"}},{"configMap":{"name":"kube-root-ca.crt","items":[{"key":"ca.crt","path":"ca.crt"}]}},{"downwardAPI":{"items":[{"path":"namespace","fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"}}]}}],"defaultMode":420}}],"containers":[{"name":"demo-crd-resource","image":"edcelvista/ubuntu24-network-tools:v13-k8s-crd","command":["sleep","infinity"],"resources":{},"volumeMounts":[{"name":"kube-api-access-5drvz","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true,"preemptionPolicy":"PreemptLowerPriority"},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}}}`)

	// Create a request to pass to the handler
	req := httptest.NewRequest(http.MethodPost, "/webhook/validating/pod", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	WebhookValidatingHandlerPOSTPod(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := "apiVersion"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestWebhookMutatingHandlerPOSTPod(t *testing.T) {
	jsonData := []byte(`{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1","request":{"uid":"cb2373d5-7933-43a4-9ca7-0be7db7c6f33","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"requestKind":{"group":"","version":"v1","kind":"Pod"},"requestResource":{"group":"","version":"v1","resource":"pods"},"name":"demo-crd-resource-76bf5f8695-sf4bt","namespace":"demo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"50211333-a14e-4c1a-9c38-e00d5b1af517","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"],"extra":{"authentication.kubernetes.io/credential-id":["JTI=7edabfc3-6128-44f4-8d8e-bc35d1ccd47d"]}},"object":{"kind":"Pod","apiVersion":"v1","metadata":{"name":"demo-crd-resource-76bf5f8695-sf4bt","generateName":"demo-crd-resource-76bf5f8695-","namespace":"demo","uid":"7df01e81-bb71-49d5-a36b-74de72836d75","creationTimestamp":"2025-04-29T04:39:35Z","labels":{"appSelector":"myapp-resource","costCenter":"test","custom-webhook.edcelvista.com/validate-always-allow":"true","pod-template-hash":"76bf5f8695","supportEmail":"administrator_edcelvista.com","tenantName":"test","mutatedLabel":"mutatedLabelValue"},"annotations":{"kubectl.kubernetes.io/restartedAt":"2025-04-29T04:39:32Z"},"ownerReferences":[{"apiVersion":"apps/v1","kind":"ReplicaSet","name":"demo-crd-resource-76bf5f8695","uid":"ee15fb57-668f-4d47-b710-29e8c37fbc7b","controller":true,"blockOwnerDeletion":true}],"managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"v1","time":"2025-04-29T04:39:35Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/restartedAt":{}},"f:generateName":{},"f:labels":{".":{},"f:appSelector":{},"f:costCenter":{},"f:custom-webhook.edcelvista.com/validate-always-allow":{},"f:pod-template-hash":{},"f:supportEmail":{},"f:tenantName":{}},"f:ownerReferences":{".":{},"k:{\"uid\":\"ee15fb57-668f-4d47-b710-29e8c37fbc7b\"}":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"demo-crd-resource\"}":{".":{},"f:command":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:enableServiceLinks":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}]},"spec":{"volumes":[{"name":"kube-api-access-5drvz","projected":{"sources":[{"serviceAccountToken":{"expirationSeconds":3607,"path":"token"}},{"configMap":{"name":"kube-root-ca.crt","items":[{"key":"ca.crt","path":"ca.crt"}]}},{"downwardAPI":{"items":[{"path":"namespace","fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"}}]}}],"defaultMode":420}}],"containers":[{"name":"demo-crd-resource","image":"edcelvista/ubuntu24-network-tools:v13-k8s-crd","command":["sleep","infinity"],"resources":{},"volumeMounts":[{"name":"kube-api-access-5drvz","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true,"preemptionPolicy":"PreemptLowerPriority"},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null,"dryRun":false,"options":{"kind":"CreateOptions","apiVersion":"meta.k8s.io/v1"}}}`)

	// Create a request to pass to the handler
	req := httptest.NewRequest(http.MethodPost, "/webhook/mutating/pod", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	WebhookMutatingHandlerPOSTPod(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := "apiVersion"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}
