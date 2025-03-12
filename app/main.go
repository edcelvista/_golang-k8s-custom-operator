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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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
var INTERVAL int = 60

func checkKubeConfig() bool {
	customKubeConfig := os.Getenv("CUSTOM_KUBE_CONFIG_PATH")
	if customKubeConfig != "" {
		KUBECONFIG = customKubeConfig
		_, err := os.Stat(KUBECONFIG)
		if !os.IsNotExist(err) {
			log.Println("Config in found [CUSTOM_KUBE_CONFIG_PATH]", KUBECONFIG)
			return true
		}
	}

	defaultHome := os.ExpandEnv("$HOME")
	_, err2 := os.Stat(defaultHome)
	if !os.IsNotExist(err2) {
		log.Println("Config in found [Default Home]", defaultHome)
		KUBECONFIG = fmt.Sprintf("%v/.kube/config", defaultHome)
		return true
	}

	log.Println("Config in not found", KUBECONFIG, defaultHome)
	return false
}

func checkRequiredEnv() {
	targetNs := os.Getenv("TARGETNAMESPACE")
	if targetNs != "" {
		TARGETNAMESPACE = targetNs
		log.Println("Env Config in found", TARGETNAMESPACE)
	}

	appName := os.Getenv("APPNAME")
	if appName != "" {
		APPNAME = appName
		log.Println("Env Config in found", APPNAME)
	}

	crdTarget := os.Getenv("CRDNAME")
	if crdTarget != "" {
		CRDNAME = crdTarget
		log.Println("Env Config in found", CRDNAME)
	}

	groupTarget := os.Getenv("CRDGROUP")
	if groupTarget != "" {
		CRDGROUP = groupTarget
		log.Println("Env Config in found", CRDGROUP)
	}

	verTarget := os.Getenv("CRDVERSION")
	if verTarget != "" {
		CRDVERSION = verTarget
		log.Println("Env Config in found", CRDVERSION)
	}

	resourceTarget := os.Getenv("CRDRESOURCE")
	if resourceTarget != "" {
		CRDRESOURCE = resourceTarget
		log.Println("Env Config in found", CRDRESOURCE)
	}

	interval := os.Getenv("INTERVAL")
	if interval != "" {
		num, err := strconv.Atoi(interval)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		INTERVAL = num
		log.Println("Env Interval in found", INTERVAL)
	}

}

func initConnection() *kubernetes.Clientset {
	// Load kubeconfig file (for running locally)
	config, err := clientcmd.BuildConfigFromFlags("", KUBECONFIG)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		log.Fatalf("Error retrieving server version: %v", err)
	}
	log.Println("Cluster Version:", version.String())

	return clientset
}

func initConnectionExt() *apiextensionsclientset.Clientset {
	// Load kubeconfig file (for running locally)
	config, err := clientcmd.BuildConfigFromFlags("", KUBECONFIG)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Create clientset
	clientset, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	return clientset
}

func initConnectionExtDynamic() *dynamic.DynamicClient {
	// Load kubeconfig file (for running locally)
	config, err := clientcmd.BuildConfigFromFlags("", KUBECONFIG)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Create clientset
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	return dynamicClient
}

func main() {
	// Create a base context (usually context.Background() or context.TODO())
	ctx := context.Background()

	// ENV Init
	checkRequiredEnv()
	isConfigExist := checkKubeConfig()
	log.Println("IsKubeConfigExist:", strconv.FormatBool(isConfigExist))

	// Init Connections
	cs := initConnection()
	cse := initConnectionExt()
	csed := initConnectionExtDynamic()

	// Store the object i.e Connections Obj in the context using context.WithValue
	ctx = context.WithValue(ctx, "clientSet", cs)
	ctx = context.WithValue(ctx, "clientSetExt", cse)
	ctx = context.WithValue(ctx, "clientSetExtDynamic", csed)

	// Check Functions
	checkNodes(ctx)
	checkCRDs(ctx)

	// Check CRD Resource
	tick := time.Tick(time.Duration(INTERVAL) * time.Second)
	for range tick {
		runOperator(ctx)
	}
}

func runOperator(ctx context.Context) {
	listCRDResources := checkCRDResources(ctx)
	for i := 0; i < len(listCRDResources); i++ {
		crdResourcePayload := listCRDResources[i]
		splitted := strings.Split(crdResourcePayload, "|")
		resourceDetail := getCRDResourceDetails(ctx, splitted[0], splitted[1])
		resourceDetail.reconcile(ctx)
	}
}

func checkNodes(ctx context.Context) {
	clientSet := ctx.Value("clientSet").(*kubernetes.Clientset)
	nodes, errNodes := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if errNodes != nil {
		log.Fatalf("Error listing pods: %v\n", errNodes)
	}

	log.Printf("NODES: %v", len(nodes.Items))
	for _, v := range nodes.Items {
		log.Println(" >", v.Name, v.Status.NodeInfo.KubeletVersion, v.Status.NodeInfo.Architecture, v.Status.NodeInfo.MachineID)
	}
}

func checkPods(ctx context.Context, listOptions metav1.ListOptions) *v1.PodList {
	clientSet := ctx.Value("clientSet").(*kubernetes.Clientset)
	pods, errPods := clientSet.CoreV1().Pods(TARGETNAMESPACE).List(ctx, listOptions)
	if errPods != nil {
		log.Fatalf("Error listing pods: %v", errPods)
	}

	log.Printf("PODS: %v\n", len(pods.Items))
	for _, v := range pods.Items {
		for _, c := range v.Status.ContainerStatuses {
			log.Printf(" > %v/%v container %v is ready: (%v)\n", v.Namespace, v.Name, c.Name, c.Ready)
		}
	}

	return pods
}

func checkDeployment(ctx context.Context, listOptions metav1.ListOptions) *appsv1.DeploymentList {
	clientSet := ctx.Value("clientSet").(*kubernetes.Clientset)
	dep, errDep := clientSet.AppsV1().Deployments(TARGETNAMESPACE).List(ctx, listOptions)
	if errDep != nil {
		log.Fatalf("Error listing Deployment: %v", errDep)
	}

	log.Printf("DEPLOYMENT: %v\n", len(dep.Items))
	for _, v := range dep.Items {
		log.Printf(" > %v/%v with ready pods: (%v)\n", v.Namespace, v.Name, v.Status.AvailableReplicas)
	}

	return dep
}

func checkCRDs(ctx context.Context) {
	clientSet := ctx.Value("clientSetExt").(*apiextensionsclientset.Clientset)

	// List CRDs
	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", CRDNAME),
	}
	crdList, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().List(ctx, listOptions)
	if err != nil {
		log.Fatalf("Error listing CRDs: %v", err)
	}

	log.Printf("CRD: %v\n", len(crdList.Items))
	for _, v := range crdList.Items {
		log.Println(" >", v.Name)
	}
}

func checkCRDResources(ctx context.Context) []string {
	clientSet := ctx.Value("clientSetExtDynamic").(*dynamic.DynamicClient)

	// Define the GVR (GroupVersionResource)
	gvr := schema.GroupVersionResource{
		Group:    CRDGROUP,
		Version:  CRDVERSION,
		Resource: CRDRESOURCE,
	}

	crdRList, err := clientSet.Resource(gvr).Namespace(TARGETNAMESPACE).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing CRD Resource: %v", err)
	}

	log.Printf("CRD RESOURCE: %v \n", len(crdRList.Items))
	resourceNames := []string{}
	for _, v := range crdRList.Items {
		resourceCollation := fmt.Sprintf("%v|%v", v.GetName(), v.GetNamespace())
		resourceNames = append(resourceNames, resourceCollation)
		log.Printf(" > %v/%v", v.GetNamespace(), v.GetName())
	}

	return resourceNames
}

func getCRDResourceDetails(ctx context.Context, name string, namespace string) MyApp {
	clientSet := ctx.Value("clientSetExtDynamic").(*dynamic.DynamicClient)

	// Define the GVR (GroupVersionResource)
	gvr := schema.GroupVersionResource{
		Group:    CRDGROUP,
		Version:  CRDVERSION,
		Resource: CRDRESOURCE,
	}

	crdResourceObj, err := clientSet.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error listing CRD Resource: %v", err)
	}

	// Manual mapping of unstructured object
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
	}

	return myApp
}

func (m *MyApp) reconcile(ctx context.Context) {
	// List pods in the "TARGETNAMESPACE" namespace
	listOptions := metav1.ListOptions{ // Use metav1.ListOptions to filter by label
		LabelSelector: fmt.Sprintf("appSelector=%v", m.spec.appSelector),
	}

	// Check Pods & Deployments
	deployments := checkDeployment(ctx, listOptions)
	pods := checkPods(ctx, listOptions)
	isReconciled := false

	if len(pods.Items) == 0 {
		if len(deployments.Items) == 0 {
			log.Printf("Reconciling Creating Deployment...")
			createDeployment(ctx, m)
			isReconciled = true
		}
	}

	for _, v := range deployments.Items {
		if m.spec.replicas != int64(v.Status.Replicas) {
			log.Printf("Reconciling Scaling Deployment: %v/%v from %v => %v", v.Namespace, v.Name, v.Status.Replicas, m.spec.replicas)
			scaleDeployment(ctx, m.spec.replicas, v.Name)
			isReconciled = true
		}

		if m.spec.image != v.Spec.Template.Spec.Containers[0].Image {
			log.Printf("Reconciling Patching Deployment: %v/%v from %v => %v", v.Namespace, v.Name, v.Spec.Template.Spec.Containers[0].Image, m.spec.image)
			patch := fmt.Sprintf(`[
				{ "op": "replace", "path": "/spec/template/spec/containers/0/image", "value": "%v" }
			]`, m.spec.image)
			updateDeployment(ctx, v.Name, patch)
			isReconciled = true
		}

		if m.spec.replicas != int64(v.Status.AvailableReplicas) {
			log.Printf("NEED MANUAL RECONCILING: POD Failing... Requires %v available replicas... current %v...", m.spec.replicas, v.Status.AvailableReplicas)
			isReconciled = true
		}
	}

	if !isReconciled {
		log.Printf("DESIRED ready pods: (%v)... NO ACTION NEEDED...", m.spec.replicas)
	}
}

func updateDeployment(ctx context.Context, deploymentName string, patchJson string) {
	clientSet := ctx.Value("clientSet").(*kubernetes.Clientset)
	// Create the Deployment
	deploymentsClient := clientSet.AppsV1().Deployments(TARGETNAMESPACE)
	ctxDep, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	patch := []byte(patchJson)

	_, err := deploymentsClient.Patch(ctxDep, deploymentName, types.JSONPatchType, patch, metav1.PatchOptions{})
	if err != nil {
		log.Fatalf("Failed to apply patch: %v", err)
	}

	// Output the result
	log.Printf("Deployment %s patched with %v\n", deploymentName, patchJson)
}

func scaleDeployment(ctx context.Context, replicas int64, deploymentName string) {
	clientSet := ctx.Value("clientSet").(*kubernetes.Clientset)
	// Create the Deployment
	deploymentsClient := clientSet.AppsV1().Deployments(TARGETNAMESPACE)
	ctxDep, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current scale
	scale, err := clientSet.AppsV1().Deployments(TARGETNAMESPACE).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get scale: %v", err)
	}

	// Update the number of replicas
	scale.Spec.Replicas = int64ToInt32Ptr(replicas) // New replica count

	_, err2 := deploymentsClient.UpdateScale(ctxDep, deploymentName, scale, metav1.UpdateOptions{
		FieldManager: "scale-controller",
	})
	if err2 != nil {
		log.Fatalf("Failed to apply scale: %v", err2)
	}

	// Output the result
	log.Printf("Deployment %s scaled to %d replicas\n", deploymentName, replicas)
}

func createDeployment(ctx context.Context, m *MyApp) {
	clientSet := ctx.Value("clientSet").(*kubernetes.Clientset)
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
	ctxDep, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := deploymentsClient.Create(ctxDep, deployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create deployment: %v", err)
	}

	log.Printf("Deployment %q created successfully\n", result.GetObjectMeta().GetName())
}
