package model

type WebhookParam struct {
	Name string
}

type WebhookResponse struct {
	WebhookParam WebhookParam
	Message      string
}

// Patch operation struct
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
type Tenant struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	ObjectMeta struct {
		CreationTimestamp string            `json:"creationTimestamp"`
		Generation        int               `json:"generation"`
		Labels            map[string]string `json:"labels"`
		Name              string            `json:"name"`
		ResourceVersion   string            `json:"resourceVersion"`
		UID               string            `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Cordoned       bool `json:"cordoned"`
		IngressOptions struct {
			HostnameCollisionScope string `json:"hostnameCollisionScope"`
		} `json:"ingressOptions"`
		LimitRanges      map[string]interface{} `json:"limitRanges"`
		NamespaceOptions struct {
			ForbiddenAnnotations map[string]interface{} `json:"forbiddenAnnotations"`
			ForbiddenLabels      map[string]interface{} `json:"forbiddenLabels"`
			Quota                int                    `json:"quota"`
		} `json:"namespaceOptions"`
		NetworkPolicies map[string]interface{} `json:"networkPolicies"`
		Owners          []struct {
			ClusterRoles []string `json:"clusterRoles"`
			Kind         string   `json:"kind"`
			Name         string   `json:"name"`
		} `json:"owners"`
		PreventDeletion bool `json:"preventDeletion"`
		ResourceQuotas  struct {
			Scope string `json:"scope"`
		} `json:"resourceQuotas"`
	} `json:"spec"`
	Status struct {
		Namespaces []string `json:"namespaces"`
		Size       int      `json:"size"`
		State      string   `json:"state"`
	} `json:"status"`
}
