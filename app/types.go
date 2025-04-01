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
}

type MyAppMetadata struct {
	name      string
	namespace string
}

type MyAppSpec struct {
	image       string
	replicas    int64
	appSelector string
}

// type GetPods struct {
// 	ctx         context.Context
// 	listOptions metav1.ListOptions
// }

type Deployments struct {
	ctx         context.Context
	listOptions metav1.ListOptions
}

type UpdateDeploymentParams struct {
	ctx       context.Context
	name      string
	patchJson string
}

type ScaleDeploymentParams struct {
	ctx      context.Context
	replicas int64
	name     string
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
	ctx       context.Context
	gvr       schema.GroupVersionResource
	name      string
	namespace string
}
