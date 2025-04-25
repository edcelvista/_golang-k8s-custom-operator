package model

type WebhookParam struct {
	Name string
}

type WebhookResponse struct {
	WebhookParam WebhookParam
	Message      string
}
