package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	authv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

var KUBECONFIG string = "~/.kube/config"
var EXCLUDENAMESPACE string = "argocd,castai-agent,custom-operator,default,krakend,kube-bench,kube-flannel,kube-node-lease,kube-public,kube-system,kyverno,monitoring,nginx-ingress"

var CRDNAME string = "sectools.k8s.edcelvista.com"
var CRDGROUP string = "k8s.edcelvista.com"
var CRDVERSION string = "v1"
var CRDRESOURCE string = "sectools"

var APPNAME string = "edcelvistadotcom-aws-tls"
var INTERVAL int = 15
var K8S_TIMEOUT int32 = 60

func checkRequiredEnv() error {
	excludeNs := os.Getenv("EXCLUDENAMESPACE")
	if excludeNs != "" {
		EXCLUDENAMESPACE = excludeNs
		log.Println("💡 Env Config in found", EXCLUDENAMESPACE)
	}

	appName := os.Getenv("APPNAME")
	if appName != "" {
		APPNAME = appName
		log.Println("💡 Env Config in found", APPNAME)
	}

	crdTarget := os.Getenv("CRDNAME")
	if crdTarget != "" {
		CRDNAME = crdTarget
		log.Println("💡 Env Config in found", CRDNAME)
	}

	groupTarget := os.Getenv("CRDGROUP")
	if groupTarget != "" {
		CRDGROUP = groupTarget
		log.Println("💡 Env Config in found", CRDGROUP)
	}

	verTarget := os.Getenv("CRDVERSION")
	if verTarget != "" {
		CRDVERSION = verTarget
		log.Println("💡 Env Config in found", CRDVERSION)
	}

	resourceTarget := os.Getenv("CRDRESOURCE")
	if resourceTarget != "" {
		CRDRESOURCE = resourceTarget
		log.Println("💡 Env Config in found", CRDRESOURCE)
	}

	interval := os.Getenv("INTERVAL")
	if interval != "" {
		num, err := strconv.Atoi(interval)
		if err != nil {
			fmt.Println("‼️ Error:", err)
			return err
		}
		INTERVAL = num
		log.Println("💡 Env Interval in found", INTERVAL)
	}
	return nil
}

func checkKubeConfig() bool {
	customKubeConfig := os.Getenv("CUSTOM_KUBE_CONFIG_PATH")
	if customKubeConfig != "" {
		KUBECONFIG = customKubeConfig
		_, err := os.Stat(KUBECONFIG)
		if !os.IsNotExist(err) {
			log.Println("💡 ⚡️ Config in found [CUSTOM_KUBE_CONFIG_PATH]", KUBECONFIG)
			return true
		}
	}

	defaultHome := os.ExpandEnv("$HOME")
	homeKubeConfig := fmt.Sprintf("%s/.kube/config", defaultHome)
	_, err := os.Stat(homeKubeConfig)
	if !os.IsNotExist(err) {
		log.Println("💡 ⚡️ Config in found [Default Home]", homeKubeConfig)
		KUBECONFIG = homeKubeConfig
		return true
	}

	log.Println("‼️ Config in not found", KUBECONFIG, homeKubeConfig)
	return false
}

func initConnection() (*kubernetes.Clientset, *apiextensionsclientset.Clientset, *dynamic.DynamicClient) {
	// Load kubeconfig file (for running locally)
	config, err := clientcmd.BuildConfigFromFlags("", KUBECONFIG)
	if err != nil {
		log.Printf("‼️ Error loading kubeconfig: %v", err)
		log.Println("💡 Defaulting to service account...")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("‼️ Error loading kubeconfig: %v", err)
		}
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("‼️ Error creating Kubernetes client: %v", err)
	}

	// Create clientset api ext
	clientsetExt, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("‼️ Error creating Kubernetes client: %v", err)
	}

	// Create clientset dynamic client
	clientSetDynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("‼️ Error creating Kubernetes client: %v", err)
	}

	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		log.Fatalf("‼️ Error retrieving server version: %v", err)
	}
	log.Println("💡 ⚡️ Cluster Version:", version.String())

	return clientset, clientsetExt, clientSetDynamic
}

func main() {
	// Create a base context (usually context.Background() or context.TODO())
	ctx := context.Background()

	// ENV Init
	err := checkRequiredEnv()
	if err != nil {
		log.Fatalf("‼️ Failed to initialize Environment Variables: %v", err)
	}

	isConfigExist := checkKubeConfig()
	if !isConfigExist {
		log.Printf("‼️ Failed to find the kube/config. isConfigExist: %v", isConfigExist)
	}
	log.Println("💡 IsKubeConfigExist:", strconv.FormatBool(isConfigExist))

	// Init Connections
	cs, cse, csed := initConnection()

	// get user context
	review, _ := cs.AuthenticationV1().SelfSubjectReviews().Create(ctx, &authv1.SelfSubjectReview{}, metav1.CreateOptions{})
	log.Printf("💡 Current User: %+v part of %+v", review.Status.UserInfo.Username, review.Status.UserInfo.Groups)

	// Store the object i.e Connections Obj in the context using context.WithValue
	ctx = context.WithValue(ctx, clientSet, cs)
	ctx = context.WithValue(ctx, clientSetExt, cse)
	ctx = context.WithValue(ctx, clientSetExtDynamic, csed)
	ctx = context.WithValue(ctx, crdListOptions, metav1.ListOptions{})
	ctx = context.WithValue(ctx, crdListOptions, metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%s", CRDNAME)})
	ctx = context.WithValue(ctx, crdResourceListOptions, metav1.ListOptions{})
	ctx = context.WithValue(ctx, crdResourceGetOptions, metav1.GetOptions{})
	// ctx = context.WithValue(ctx, "podListOptions", metav1.ListOptions{})

	// Check Functions
	checkNodes(ctx)
	checkCrds(ctx)

	// Check CRD Resource
	tick := time.Tick(time.Duration(INTERVAL) * time.Second)
	operator := Operator{
		ctx: ctx,
		gvr: schema.GroupVersionResource{ // Define the GVR (GroupVersionResource)
			Group:    CRDGROUP,
			Version:  CRDVERSION,
			Resource: CRDRESOURCE,
		},
	}

	for range tick {
		log.Printf("🪝 Ticking...")
		operator.Run()
	}
}

func checkNodes(ctx context.Context) {
	clientSet := ctx.Value(clientSet).(*kubernetes.Clientset)

	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("‼️ Error listing nodes: %v\n", err)
	}

	log.Printf("🖥️  NODES: %+v", len(nodes.Items))
	for _, v := range nodes.Items {
		log.Println("⚡️ ", v.Name, v.Status.NodeInfo.KubeletVersion, v.Status.NodeInfo.Architecture, v.Status.NodeInfo.MachineID)
	}
}

func checkCrds(ctx context.Context) {
	clientSet := ctx.Value(clientSetExt).(*apiextensionsclientset.Clientset)
	crdListOptions := ctx.Value(crdListOptions).(metav1.ListOptions)

	crds, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().List(ctx, crdListOptions)
	if err != nil {
		log.Fatalf("‼️ Error listing CRDs: %v", err)
	}

	log.Printf("📝 CRD: %v\n", len(crds.Items))
	for _, v := range crds.Items {
		log.Println("⚡️", v.Name)
	}
}

func (m *Operator) Run() {
	cr := CRDResources{
		ctx: m.ctx,
		gvr: m.gvr,
	}
	crdResources := cr.get()
	for i := 0; i < len(crdResources); i++ {
		crdResourceName := crdResources[i]
		crdResource := CRDResourceDetails{
			ctx:  m.ctx,
			gvr:  m.gvr,
			name: crdResourceName,
		}
		crdResourceDetails := crdResource.get()
		crdResourceDetails.reconcile()
	}

	if len(crdResources) == 0 {
		log.Println("⚠️ No CRD Resource Found! Operator will not trigger.")
	}
}

func (m *CRDResources) get() []string {
	clientSet := m.ctx.Value(clientSetExtDynamic).(*dynamic.DynamicClient)
	crdResourceListOptions := m.ctx.Value(crdResourceListOptions).(metav1.ListOptions)

	crdResources, err := clientSet.Resource(m.gvr).Namespace("").List(m.ctx, crdResourceListOptions)
	if err != nil {
		log.Fatalf("‼️ Error listing CRD Resource: %v", err)
	}

	log.Printf("🗒️ CRD RESOURCE: %v \n", len(crdResources.Items))
	resourceNames := []string{}
	for _, v := range crdResources.Items {
		resourceNames = append(resourceNames, v.GetName())
		log.Printf("⚡️ %v", v.GetName())
	}

	return resourceNames
}

func (m *CRDResourceDetails) get() *MyApp {
	clientSet := m.ctx.Value(clientSetExtDynamic).(*dynamic.DynamicClient)
	crdResourceGetOptions := m.ctx.Value(crdResourceGetOptions).(metav1.GetOptions)

	crdResourceObj, err := clientSet.Resource(m.gvr).Namespace("").Get(m.ctx, m.name, crdResourceGetOptions)
	if err != nil {
		log.Fatalf("‼️ Error listing CRD Resource: %v", err)
	}

	// Manual mapping of unstructured object, assert type and build struct data
	metadata := crdResourceObj.Object["metadata"].(map[string]interface{})
	spec := crdResourceObj.Object["spec"].(map[string]interface{})

	secret := spec["secret"].(map[string]interface{})
	secretName := secret["name"].(string)
	secretType := secret["type"].(string)
	data := secret["data"].(map[string]interface{})
	addProps := data["additionalProperties"].(map[string]interface{})
	tlsCrt := addProps["tls.crt"].(string)
	tlsKey := addProps["tls.key"].(string)

	myApp := MyApp{
		metadata: MyAppMetadata{
			name: metadata["name"].(string),
		},
		spec: MyAppSpec{
			SecretProperties{
				SecretData{
					AdditionalProperties{
						TLSCrt: tlsCrt,
						TLSKey: tlsKey,
						Name:   secretName,
						Type:   secretType,
					},
				},
			},
		},
		ctx: m.ctx,
	}

	return &myApp
}

func (m *SecretTarget) get() (*corev1.Secret, error) {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	getOptions := metav1.GetOptions{}

	secret, err := clientSet.CoreV1().Secrets(m.namespace).Get(m.ctx, m.name, getOptions)
	if err != nil {
		return secret, err
	} else {
		log.Printf("⚡️ %v/%v", secret.Namespace, secret.Name)
		log.Printf("🛠️ SECRET: %v/%v\n", secret.Namespace, secret.Name)
	}

	return secret, nil
}

func (m *Namespaces) get() (*corev1.NamespaceList, error) {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	// Get all namespaces
	ns, err := clientSet.CoreV1().Namespaces().List(m.ctx, m.listOptions)
	if err != nil {
		return ns, err
	}

	return ns, nil
}

func (m *MyApp) reconcile() {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	listOptions := metav1.ListOptions{}

	n := Namespaces{
		ctx:         m.ctx,
		listOptions: listOptions,
	}

	ns, err := n.get()
	if err != nil {
		log.Fatalf("‼️ Error listing namespace: %v", err)
	}

	isChanged := false
	for _, ns := range ns.Items {
		isChanged := false
		if strings.Contains(EXCLUDENAMESPACE, ns.Name) {
			log.Printf("💡 Skipping namespace (%v) ...", ns.Name)
			continue
		} else {
			s := SecretTarget{
				ctx:         m.ctx,
				listOptions: listOptions,
				namespace:   ns.Name,
				name:        m.spec.secret.data.additionalProperties.Name,
			}

			decodedBytesCrt, err := base64.StdEncoding.DecodeString(m.spec.secret.data.additionalProperties.TLSCrt)
			if err != nil {
				log.Fatalf("‼️ Error decoding secret Cert: %v", err)
			}
			decodedBytesKey, err := base64.StdEncoding.DecodeString(m.spec.secret.data.additionalProperties.TLSKey)
			if err != nil {
				log.Fatalf("‼️ Error decoding secret Key: %v", err)
			}

			secretTarget, err := s.get()

			if err != nil {
				log.Printf("⚠️ ⚡️ Reconciling Creating Secret (%v) in (%v)...", s.name, s.namespace)

				secretDetail := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: APPNAME,
					},
					// Immutable: BoolPtr(true),
					Type: corev1.SecretTypeTLS,
					Data: map[string][]byte{
						"tls.crt": decodedBytesCrt,
						"tls.key": decodedBytesKey,
					},
				}

				createOptions := metav1.CreateOptions{}
				createSecret, err := clientSet.CoreV1().Secrets(s.namespace).Create(m.ctx, secretDetail, createOptions)
				if err != nil {
					log.Printf("‼️ Error creating secret: %v", err)
				}
				log.Printf("⚡️ Created Secret (%v) in (%v)...", createSecret.Name, createSecret.Namespace)
				isChanged = true
			} else {
				existingSecretType := string(secretTarget.Type)
				existingSecretTlsCrt := string(secretTarget.Data["tls.crt"])
				existingSecretTlsKey := string(secretTarget.Data["tls.key"])

				change := ""
				if base64.StdEncoding.EncodeToString([]byte(existingSecretTlsCrt)) != m.spec.secret.data.additionalProperties.TLSCrt {
					change = fmt.Sprintf("%v %v | %v vs %v", change, "Cert", "omitted", "omitted")
					isChanged = true
				}

				if base64.StdEncoding.EncodeToString([]byte(existingSecretTlsKey)) != m.spec.secret.data.additionalProperties.TLSKey {
					change = fmt.Sprintf("%v %v | %v vs %v", change, "Key", "omitted", "omitted")
					isChanged = true
				}

				if existingSecretType != m.spec.secret.data.additionalProperties.Type {
					change = fmt.Sprintf("%v %v | %v vs %v", change, "Type", existingSecretType, m.spec.secret.data.additionalProperties.Type)
					isChanged = true
				}

				if isChanged {
					secretClient := clientSet.CoreV1().Secrets(s.namespace)
					secret, err := secretClient.Get(m.ctx, s.name, metav1.GetOptions{})
					if err != nil {
						log.Fatalf("‼️ Error getting secret: %v", err)
					}
					// Update the data

					secret.Data["tls.crt"] = decodedBytesCrt
					secret.Data["tls.key"] = decodedBytesKey
					secret.Type = corev1.SecretTypeTLS

					// Apply the update
					updatedSecret, err := secretClient.Update(m.ctx, secret, metav1.UpdateOptions{})
					if err != nil {
						log.Fatalf("‼️ Error Patching secret: %v", err)
					}

					log.Printf("⚠️ ⚡️ Reconciling Updating Secret (%v) in (%v) changed in %v...", updatedSecret.Name, updatedSecret.Namespace, change)
				} else {
					log.Printf("💡 Secret %v/%v at latest no action required ...", secretTarget.Namespace, secretTarget.Name)
				}
			}
		}
	}

	if !isChanged {
		log.Printf("✅ CRD DESIRED ... NO ACTION NEEDED...")
	}
}
