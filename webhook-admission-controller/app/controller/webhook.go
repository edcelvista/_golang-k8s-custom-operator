package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	Model "_gorestapi-k8s/model"

	"github.com/gorilla/mux"
	admissionv1 "k8s.io/api/admission/v1"
)

func WebhookHandlerGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqParams := mux.Vars(r)

	res := Model.WebhookResponse{
		WebhookParam: Model.WebhookParam{Name: reqParams["name"]},
		Message:      fmt.Sprintf("💡 Webhook for %v", reqParams["name"]),
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := &admissionv1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
	}

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := &admissionv1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
	}

	// Build response
	admissionReviewResp.TypeMeta = admissionReviewReq.TypeMeta
	admissionReviewResp.Response = res

	json.NewEncoder(w).Encode(admissionReviewResp)
}
