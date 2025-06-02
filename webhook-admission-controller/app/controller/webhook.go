package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	Lib "_gorestapi-k8s/lib"
	Model "_gorestapi-k8s/model"

	"github.com/gorilla/mux"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WebhookHandlerGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	log.Printf("💡 Received Webhook for %v", reqParams["name"])

	res := Model.WebhookResponse{
		WebhookParam: Model.WebhookParam{Name: reqParams["name"]},
		Message:      fmt.Sprintf("💡 Received Webhook for %v", reqParams["name"]),
	}

	json.NewEncoder(w).Encode(res)
}

func WebhookValidatingHandlerPOSTPod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	// Parse JSON from request body
	var admissionReviewReq admissionv1.AdmissionReview
	var admissionReviewResp admissionv1.AdmissionReview

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) //you want to clone r.Body in Go (so you can read it multiple times).
	Lib.Debug(fmt.Sprintf("[VAL] [REQ] %+v %+v %+v %+v", string(bodyBytes), r.Header, r.Host, r.URL))

	err := json.NewDecoder(r.Body).Decode(&admissionReviewReq)
	if err != nil {
		log.Printf("‼️ Error Validating Webhook %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if admissionReviewReq.Request == nil {
		log.Println("‼️ AdmissionReview.Request is nil")
		http.Error(w, "Invalid AdmissionReview", http.StatusBadRequest)
		return
	}

	var pod corev1.Pod
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod); err != nil {
		log.Println("‼️ Error Validating Webhook could not unmarshal raw pod object")
		http.Error(w, "could not unmarshal raw pod object", http.StatusBadRequest)
		return
	}

	log.Printf("💡 Received Validating Webhook %v event object source %v/%v", reqParams["name"], pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
	res := &admissionv1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
	}

	requiredLabels := strings.Split(os.Getenv("VALIDATE_LABEL"), ",")
	var warnings []string

	if pod.Labels == nil {
		res.Allowed = false
		warnings = append(warnings, fmt.Sprintf("%v missing metadata label", os.Getenv("VALIDATE_LABEL")))
	} else {
		for i := 0; i < len(requiredLabels); i++ {
			if pod.Labels[requiredLabels[i]] == "" {
				res.Allowed = false
				warnings = append(warnings, fmt.Sprintf("%v missing metadata label", requiredLabels[i]))
			}
		}
	}

	message := strings.Join(warnings, " | ")
	res.Result = &metav1.Status{
		Message: message,
	}

	log.Printf("⚡️ Validating Webhook Message [%v/%v]: %v", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, Lib.DefaultIfEmpty(message, "no-message"))

	// Build response
	admissionReviewResp.TypeMeta = admissionReviewReq.TypeMeta
	admissionReviewResp.Response = res

	jsonData, _ := json.Marshal(admissionReviewResp)
	Lib.Debug(fmt.Sprintf("[VAL] [RES] %+v", string(jsonData)))
	json.NewEncoder(w).Encode(admissionReviewResp)
}

func WebhookMutatingHandlerPOSTPod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	// Parse JSON from request body
	var admissionReviewReq admissionv1.AdmissionReview
	var admissionReviewResp admissionv1.AdmissionReview

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) //you want to clone r.Body in Go (so you can read it multiple times).
	Lib.Debug(fmt.Sprintf("[MUT] [REQ] %+v %+v %+v %+v", string(bodyBytes), r.Header, r.Host, r.URL))

	err := json.NewDecoder(r.Body).Decode(&admissionReviewReq)
	if err != nil {
		log.Printf("‼️ Error Mutating Webhook %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if admissionReviewReq.Request == nil {
		log.Println("‼️ AdmissionReview.Request is nil")
		http.Error(w, "Invalid AdmissionReview", http.StatusBadRequest)
		return
	}

	var pod corev1.Pod
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod); err != nil {
		log.Println("‼️ Error Mutating Webhook could not unmarshal raw pod object")
		http.Error(w, "could not unmarshal raw pod object", http.StatusBadRequest)
		return
	}

	log.Printf("💡 Received Mutating Webhook %v event object source %v/%v", reqParams["name"], pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
	res := &admissionv1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
	}

	// Build patch operations
	patchesJson := os.Getenv("MUTATE_PATCH")
	if patchesJson == "" {
		log.Println("‼️ Error Mutating Webhook missing patch object parameter")
		http.Error(w, "missing patch object parameter", http.StatusInternalServerError)
		return
	}

	res.Patch = []byte(patchesJson)
	res.PatchType = func() *admissionv1.PatchType {
		pt := admissionv1.PatchTypeJSONPatch
		return &pt
	}()

	log.Printf("⚡️ Mutating Webhook Message [%v/%v]: %v", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, patchesJson)

	// Build response
	admissionReviewResp.TypeMeta = admissionReviewReq.TypeMeta
	admissionReviewResp.Response = res

	jsonData, _ := json.Marshal(admissionReviewResp)
	Lib.Debug(fmt.Sprintf("[MUT] [RES] %+v", string(jsonData)))
	json.NewEncoder(w).Encode(admissionReviewResp)
}

func WebhookValidatingHandlerPOSTTenant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	// Parse JSON from request body
	var admissionReviewReq admissionv1.AdmissionReview
	var admissionReviewResp admissionv1.AdmissionReview

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) //you want to clone r.Body in Go (so you can read it multiple times).
	Lib.Debug(fmt.Sprintf("[VAL] [REQ] %+v %+v %+v %+v", string(bodyBytes), r.Header, r.Host, r.URL))

	err := json.NewDecoder(r.Body).Decode(&admissionReviewReq)
	if err != nil {
		log.Printf("‼️ Error Validating Webhook %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if admissionReviewReq.Request == nil {
		log.Println("‼️ AdmissionReview.Request is nil")
		http.Error(w, "Invalid AdmissionReview", http.StatusBadRequest)
		return
	}

	var tenant Model.Tenant
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &tenant); err != nil {
		log.Println("‼️ Error Validating Webhook could not unmarshal raw tenant object")
		http.Error(w, "could not unmarshal raw tenant object", http.StatusBadRequest)
		return
	}

	var tenantOld Model.Tenant
	if err := json.Unmarshal(admissionReviewReq.Request.OldObject.Raw, &tenantOld); err != nil {
		log.Println("‼️ Error Validating Webhook could not unmarshal raw tenant object")
		http.Error(w, "could not unmarshal raw tenant object", http.StatusBadRequest)
		return
	}

	log.Printf("💡 Received Validating Webhook %v event object source %v", reqParams["name"], tenant.ObjectMeta.Name)
	res := &admissionv1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
	}

	var warnings []string

	objString, _ := json.Marshal(tenant.Spec.ResourceQuotas)
	oldObjString, _ := json.Marshal(tenantOld.Spec.ResourceQuotas)

	if string(objString) != string(oldObjString) {
		res.Allowed = true
		warnings = append(warnings, "Detected Resource Quota Changes...")
	} else {
		res.Allowed = true
	}

	log.Printf("💡 OLD %v |<>| NEW %v", string(oldObjString), string(objString))

	message := strings.Join(warnings, " | ")
	res.Result = &metav1.Status{
		Message: message,
	}

	log.Printf("⚡️ Validating Webhook Message [%v]: %v", tenant.ObjectMeta.Name, Lib.DefaultIfEmpty(message, "no-message"))

	// Build response
	admissionReviewResp.TypeMeta = admissionReviewReq.TypeMeta
	admissionReviewResp.Response = res

	jsonData, _ := json.Marshal(admissionReviewResp)
	Lib.Debug(fmt.Sprintf("[VAL] [RES] %+v", string(jsonData)))
	json.NewEncoder(w).Encode(admissionReviewResp)
}
