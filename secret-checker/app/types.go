package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ctxKeyClientSet string
type ctxKeyClientSetExt string
type ctxKeyClientSetExtDynamic string
type ctxKeyCodeListOptions string
type ctxKeyCrdListOptions string
type ctxKeyCrdResourceListOptions string
type ctxKeyCrdResourceGetOptions string

const (
	clientSet              ctxKeyClientSet              = "ctxKeyClientSet"
	clientSetExt           ctxKeyClientSetExt           = "ctxKeyClientSetExt"
	clientSetExtDynamic    ctxKeyClientSetExtDynamic    = "ctxKeyClientSetExtDynamic"
	codeListOptions        ctxKeyCodeListOptions        = "ctxKeyCodeListOptions"
	crdListOptions         ctxKeyCrdListOptions         = "ctxKeyCrdListOptions"
	crdResourceListOptions ctxKeyCrdResourceListOptions = "ctxKeyCrdResourceListOptions"
	crdResourceGetOptions  ctxKeyCrdResourceGetOptions  = "ctxKeyCrdResourceGetOptions"
)

type MyApp struct {
	metadata MyAppMetadata
	spec     MyAppSpec
	ctx      context.Context
}

type MyAppMetadata struct {
	name string
}

type MyAppSpec struct {
	secret SecretProperties
}

type SecretProperties struct {
	data SecretData
}

type SecretData struct {
	additionalProperties AdditionalProperties
}

type AdditionalProperties struct {
	TLSCrt string `json:"tls.crt"`
	TLSKey string `json:"tls.key"`
	Name   string
	Type   string
}

type SecretTarget struct {
	ctx         context.Context
	listOptions metav1.ListOptions
	namespace   string
	name        string
}

type Operator struct {
	ctx context.Context
	gvr schema.GroupVersionResource
}

type CRDResources struct {
	ctx context.Context
	gvr schema.GroupVersionResource
}

type CRDResourceDetails struct {
	ctx  context.Context
	gvr  schema.GroupVersionResource
	name string
}
