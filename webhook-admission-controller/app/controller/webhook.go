package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

func WebhookValidatingHandlerPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	log.Printf("💡 Received Validating Webhook %v", reqParams["name"])

	// Parse JSON from request body
	var admissionReviewReq admissionv1.AdmissionReview
	var admissionReviewResp admissionv1.AdmissionReview

	err := json.NewDecoder(r.Body).Decode(&admissionReviewReq)
	if err != nil {
		log.Printf("‼️ Error Validating Webhook %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pod corev1.Pod
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod); err != nil {
		log.Println("‼️ Error Validating Webhook could not unmarshal raw pod object")
		http.Error(w, "could not unmarshal raw pod object", http.StatusBadRequest)
		return
	}

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
	log.Printf("💡 ⚡️ Validating Webhook Message: %v", message)

	// Build response
	admissionReviewResp.TypeMeta = admissionReviewReq.TypeMeta
	admissionReviewResp.Response = res

	json.NewEncoder(w).Encode(admissionReviewResp)
}

func WebhookMutatingHandlerPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	log.Printf("💡 Received Mutating Webhook %v", reqParams["name"])

	// Parse JSON from request body
	var admissionReviewReq admissionv1.AdmissionReview
	var admissionReviewResp admissionv1.AdmissionReview

	err := json.NewDecoder(r.Body).Decode(&admissionReviewReq)
	if err != nil {
		log.Printf("‼️ Error Mutating Webhook %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pod corev1.Pod
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod); err != nil {
		log.Println("‼️ Error Mutating Webhook could not unmarshal raw pod object")
		http.Error(w, "could not unmarshal raw pod object", http.StatusBadRequest)
		return
	}

	res := &admissionv1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
	}

	// Check if the label exists
	var patchBytes []byte
	if pod.Labels == nil || pod.Labels["mylabel"] == "" {
		// Build patch operations
		patches := []Model.PatchOperation{
			{
				Op:    os.Getenv("PATCH_OP"),
				Path:  os.Getenv("PATCH_TARGET"), //"/metadata/labels/mylabel",
				Value: os.Getenv("PATCH_VALUE"),
			},
		}

		patchBytes, err = json.Marshal(patches)
		if err != nil {
			log.Println("‼️ Error Mutating Webhook could not marshal patch")
			http.Error(w, "could not marshal patch", http.StatusInternalServerError)
			return
		}

		res.Patch = patchBytes
		res.PatchType = func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}()

		log.Printf("💡 ⚡️ Mutating Webhook Message: %v", patches)
	}

	// Build response
	admissionReviewResp.TypeMeta = admissionReviewReq.TypeMeta
	admissionReviewResp.Response = res

	json.NewEncoder(w).Encode(admissionReviewResp)
}
