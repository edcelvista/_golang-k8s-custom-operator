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
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	ObjectMeta TMetadata `json:"metadata"`
	Spec       TSpec     `json:"spec"`
	Status     TStatus   `json:"status"`
}

type TMetadata struct {
	Name              string                   `json:"name"`
	Labels            map[string]string        `json:"labels"`
	UID               string                   `json:"uid"`
	ResourceVersion   string                   `json:"resourceVersion"`
	CreationTimestamp string                   `json:"creationTimestamp"`
	Generation        int64                    `json:"generation"`
	ManagedFields     []map[string]interface{} `json:"managedFields"` // You can make this a full struct too
}

type TSpec struct {
	Cordoned        bool              `json:"cordoned"`
	IngressOptions  TIngressOptions   `json:"ingressOptions"`
	LimitRanges     map[string]string `json:"limitRanges"` // Or custom struct if needed
	NamespaceOpts   TNamespaceOptions `json:"namespaceOptions"`
	NetworkPolicies map[string]string `json:"networkPolicies"` // Or custom struct if needed
	Owners          []TOwner          `json:"owners"`
	PreventDeletion bool              `json:"preventDeletion"`
	ResourceQuotas  TResourceQuotas   `json:"resourceQuotas"`
}

type TIngressOptions struct {
	HostnameCollisionScope string `json:"hostnameCollisionScope"`
}

type TNamespaceOptions struct {
	ForbiddenAnnotations map[string]string `json:"forbiddenAnnotations"`
	ForbiddenLabels      map[string]string `json:"forbiddenLabels"`
	Quota                int               `json:"quota"`
}

type TOwner struct {
	ClusterRoles []string `json:"clusterRoles"`
	Kind         string   `json:"kind"`
	Name         string   `json:"name"`
}

type TResourceQuotas struct {
	Items []TQuotaItem `json:"items"`
	Scope string       `json:"scope"`
}

type TQuotaItem struct {
	Hard map[string]string `json:"hard"`
}

type TStatus struct {
	Namespaces []string `json:"namespaces"`
	Size       int      `json:"size"`
	State      string   `json:"state"`
}
