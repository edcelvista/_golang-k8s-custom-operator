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
