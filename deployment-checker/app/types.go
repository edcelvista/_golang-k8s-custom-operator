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
type KubeConfig struct {
	ApiVersion     string      `json: "apiVersion"`
	Kind           string      `json: "kind"`
	Clusters       []Cluster   `json: "clusters"`
	Contexts       []Context   `json: "contexts"`
	CurrentContext string      `json: "current-context"`
	Preferences    interface{} `json: "preferences"`
	Users          []User      `json:"users"`
}

type Cluster struct {
	Name    string
	Cluster ClusterDetails
}

type ClusterDetails struct {
	Server                   string `json: "server"`
	CertificateAuthorityData string `json: "certificate-authority-data"`
}

type Context struct {
	Name    string         `json: "name"`
	Context ContextDetails `json: "context"`
}

type ContextDetails struct {
	Cluster string `json: "Cluster"`
	User    string `json: "User"`
}

type User struct {
	Name string      `json: "name"`
	User UserDetails `json: "user"`
}

type UserDetails struct {
	Token string `json: "token"`
}
