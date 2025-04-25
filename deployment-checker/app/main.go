package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	authv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

var KUBECONFIG string = "~/.kube/config"
var TARGETNAMESPACE string = "demo"

var CRDNAME string = "myapps.k8s.edcelvista.com"
var CRDGROUP string = "k8s.edcelvista.com"
var CRDVERSION string = "v1"
var CRDRESOURCE string = "myapps"

var APPNAME string = "demo-custom-resource"
var INTERVAL int = 15
var K8S_TIMEOUT int32 = 60

func checkRequiredEnv() error {
	targetNs := os.Getenv("TARGETNAMESPACE")
	if targetNs != "" {
		TARGETNAMESPACE = targetNs
		log.Println("💡 Env Config in found", TARGETNAMESPACE)
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
		log.Fatalf("‼️ Error loading kubeconfig: %v", err)
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
	log.Println("IsKubeConfigExist:", strconv.FormatBool(isConfigExist))

	// Init Connections
	cs, cse, csed := initConnection()

	// get user context
	review, _ := cs.AuthenticationV1().SelfSubjectReviews().Create(ctx, &authv1.SelfSubjectReview{}, metav1.CreateOptions{})
	log.Printf("💡 Current User: %v part of %v", review.Status.UserInfo.Username, review.Status.UserInfo.Groups)

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

	log.Printf("🖥️  NODES: %v", len(nodes.Items))
	for _, v := range nodes.Items {
		log.Println("⚡️ ", v.Name, v.Status.NodeInfo.KubeletVersion, v.Status.NodeInfo.Architecture, v.Status.NodeInfo.MachineID, v.Status.Conditions[4])
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
		crdResourcePayload := crdResources[i]
		splitted := strings.Split(crdResourcePayload, "|") // name|namespace
		crdResource := CRDResourceDetails{
			ctx:       m.ctx,
			gvr:       m.gvr,
			name:      splitted[0],
			namespace: splitted[1],
		}
		crdResourceDetails := crdResource.get()
		crdResourceDetails.reconcile()
	}

	if len(crdResources) == 0 {
		log.Println("⚠️ No CRD Found")
	}
}

// func (m *GetPods) checkPods() *v1.PodList {
// 	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)
// 	podListOptions := m.ctx.Value("podListOptions").(metav1.ListOptions)

// 	pods, err := clientSet.CoreV1().Pods(TARGETNAMESPACE).List(m.ctx, podListOptions)
// 	if err != nil {
// 		log.Fatalf("‼️ Error listing pods: %v", err)
// 	}

// 	log.Printf("⚙️ PODS: %v\n", len(pods.Items))
// 	for _, v := range pods.Items {
// 		for _, c := range v.Status.ContainerStatuses {
// 			log.Printf(" ⚡️ %v/%v container %v is ready: (%v)\n", v.Namespace, v.Name, c.Name, c.Ready)
// 		}
// 	}

// 	return pods
// }

func (m *CRDResources) get() []string {
	clientSet := m.ctx.Value(clientSetExtDynamic).(*dynamic.DynamicClient)
	crdResourceListOptions := m.ctx.Value(crdResourceListOptions).(metav1.ListOptions)

	crdResources, err := clientSet.Resource(m.gvr).Namespace(TARGETNAMESPACE).List(m.ctx, crdResourceListOptions)
	if err != nil {
		log.Fatalf("‼️ Error listing CRD Resource: %v", err)
	}

	log.Printf("🗒️ CRD RESOURCE: %v \n", len(crdResources.Items))
	resourceNames := []string{}
	for _, v := range crdResources.Items {
		resourceCollation := fmt.Sprintf("%v|%v", v.GetName(), v.GetNamespace())
		resourceNames = append(resourceNames, resourceCollation)
		log.Printf("⚡️ %v/%v", v.GetNamespace(), v.GetName())
	}

	return resourceNames
}

func (m *CRDResourceDetails) get() *MyApp {
	clientSet := m.ctx.Value(clientSetExtDynamic).(*dynamic.DynamicClient)
	crdResourceGetOptions := m.ctx.Value(crdResourceGetOptions).(metav1.GetOptions)

	crdResourceObj, err := clientSet.Resource(m.gvr).Namespace(m.namespace).Get(m.ctx, m.name, crdResourceGetOptions)
	if err != nil {
		log.Fatalf("‼️ Error listing CRD Resource: %v", err)
	}

	// Manual mapping of unstructured object, assert type and build struct data
	metadata := crdResourceObj.Object["metadata"].(map[string]interface{})
	spec := crdResourceObj.Object["spec"].(map[string]interface{})
	myApp := MyApp{
		metadata: MyAppMetadata{
			name:      metadata["name"].(string),
			namespace: metadata["namespace"].(string),
		},
		spec: MyAppSpec{
			image:       spec["image"].(string),
			replicas:    spec["replicas"].(int64),
			appSelector: spec["appSelector"].(string),
		},
		ctx: m.ctx,
	}

	return &myApp
}

func (m *Deployments) get() *appsv1.DeploymentList {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	deployments, err := clientSet.AppsV1().Deployments(TARGETNAMESPACE).List(m.ctx, m.listOptions)
	if err != nil {
		log.Fatalf("‼️ Error listing Deployment: %v", err)
	}

	log.Printf("⚙️ DEPLOYMENT: %v\n", len(deployments.Items))
	for _, v := range deployments.Items {
		log.Printf("⚡️ %v/%v", v.Namespace, v.Name)
		log.Printf("🛠️ PODS: Replicas (%v) | Available (%v) | Ready (%v) | Unavailable (%v)\n", v.Status.Replicas, v.Status.AvailableReplicas, v.Status.ReadyReplicas, v.Status.UnavailableReplicas)
	}

	return deployments
}

func (m *MyApp) reconcile() {
	// List pods in the "TARGETNAMESPACE" namespace
	listOptions := metav1.ListOptions{ // Use metav1.ListOptions to filter by label
		LabelSelector: fmt.Sprintf("appSelector=%v", m.spec.appSelector),
	}

	// Check Pods & Deployments
	deps := Deployments{
		ctx:         m.ctx,
		listOptions: listOptions,
	}
	deployments := deps.get()
	// pods := checkPods(ctx, listOptions)

	if len(deployments.Items) == 0 {
		log.Printf("⚠️ ⚡️ Reconciling Creating Deployment...")
		m.createDeployment()
		return
	}

	isChanged := false
	for _, v := range deployments.Items {
		if m.spec.replicas != int64(v.Status.Replicas) {
			log.Printf("⚠️ ⚡️ Reconciling Scaling Deployment: %v/%v from %v => %v", v.Namespace, v.Name, v.Status.Replicas, m.spec.replicas)
			deploymentForScale := ScaleDeploymentParams{
				ctx:      m.ctx,
				replicas: m.spec.replicas,
				name:     v.Name,
			}
			deploymentForScale.scaleDeployment()
			isChanged = true
		}

		if m.spec.image != v.Spec.Template.Spec.Containers[0].Image {
			log.Printf("⚠️ ⚡️ Reconciling Patching Deployment: %v/%v from %v => %v", v.Namespace, v.Name, v.Spec.Template.Spec.Containers[0].Image, m.spec.image)
			patch := fmt.Sprintf(`[{ "op": "replace", "path": "/spec/template/spec/containers/0/image", "value": "%v" }]`, m.spec.image)
			deploymentForUpdate := UpdateDeploymentParams{
				ctx:       m.ctx,
				name:      v.Name,
				patchJson: patch,
			}
			deploymentForUpdate.updateDeployment()
			isChanged = true
		}
	}

	if !isChanged {
		log.Printf("✅ CRD DESIRED (%v)... NO ACTION NEEDED...", m.spec.replicas)
	}
}

func (m *MyApp) createDeployment() {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	// Define the Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      APPNAME,
			Namespace: TARGETNAMESPACE,
			Labels:    map[string]string{"appSelector": m.spec.appSelector},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int64ToInt32PtrP(m.spec.replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"appSelector": m.spec.appSelector},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"appSelector": m.spec.appSelector},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    APPNAME,
							Image:   m.spec.image,
							Command: []string{"sleep", "infinity"},
						},
					},
				},
			},
		},
	}

	// Create the Deployment
	deploymentsClient := clientSet.AppsV1().Deployments(TARGETNAMESPACE)
	ctxDep, cancel := context.WithTimeout(m.ctx, time.Duration(K8S_TIMEOUT)*time.Second)
	defer cancel()

	result, err := deploymentsClient.Create(ctxDep, deployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("‼️ Failed to create deployment: %v", err)
	}

	log.Printf("⚙️ Deployment %q created successfully...\n", result.GetObjectMeta().GetName())
}

func (m *UpdateDeploymentParams) updateDeployment() error {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	// Create the Deployment
	deploymentsClient := clientSet.AppsV1().Deployments(TARGETNAMESPACE)
	ctxDep, cancel := context.WithTimeout(m.ctx, time.Duration(K8S_TIMEOUT)*time.Second)
	defer cancel()

	patch := []byte(m.patchJson)

	_, err := deploymentsClient.Patch(ctxDep, m.name, types.JSONPatchType, patch, metav1.PatchOptions{})
	if err != nil {
		log.Fatalf("‼️ Failed to apply patch: %v", err)
		return err
	}

	log.Printf("⚙️ Deployment %s patched with %v\n", m.name, m.patchJson)

	return nil
}

func (m *ScaleDeploymentParams) scaleDeployment() error {
	clientSet := m.ctx.Value(clientSet).(*kubernetes.Clientset)

	// Create the Deployment
	deploymentsClient := clientSet.AppsV1().Deployments(TARGETNAMESPACE)
	ctxDep, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	// Get current scale
	scale, err := clientSet.AppsV1().Deployments(TARGETNAMESPACE).GetScale(ctxDep, m.name, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("‼️ Failed to get scale: %v", err)
		return err
	}

	// Update the number of replicas
	scale.Spec.Replicas = int64ToInt32Ptr(m.replicas) // New replica count

	_, err = deploymentsClient.UpdateScale(ctxDep, m.name, scale, metav1.UpdateOptions{
		FieldManager: "scale-controller",
	})
	if err != nil {
		log.Fatalf("‼️ Failed to apply scale: %v", err)
		return err
	}

	log.Printf("⚙️ Deployment %s scaled to %d replicas\n", m.name, m.replicas)

	return nil
}
