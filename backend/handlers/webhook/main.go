package webhook

import (
	// "crypto/hmac"
	// "crypto/sha256"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	subscriptionDBHelper "github.com/harshgupta9473/webhookDelivery/helpers/subscriptions"
	webhookHelper "github.com/harshgupta9473/webhookDelivery/helpers/webhook"
	"github.com/harshgupta9473/webhookDelivery/models"
	"github.com/harshgupta9473/webhookDelivery/redis"
	"github.com/harshgupta9473/webhookDelivery/utils"
	"github.com/harshgupta9473/webhookDelivery/workers/queue"
)
func IngestWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subscriptionIDstr := vars["subscriptionID"]

	subscriptionID, err := uuid.Parse(subscriptionIDstr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid subscription_id",
			Data:    nil,
		})
		return
	}

	rawPayload, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Unable to read request body",
			Data:    nil,
		})
		return
	}
	defer r.Body.Close()

	subscription, err := subscriptionDBHelper.GetSubscriptionDetailsUsingDBorCache(subscriptionID)
	if err != nil {
		utils.WriteJson(w, http.StatusNotFound, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Subscription not found: %v", err),
			Data:    nil,
		})
		return
	}

	eventType := r.Header.Get("X-event-type")
	if !isValidEventType(subscription.EventTypes, eventType) {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid event type: %v", eventType),
			Data:    nil,
		})
		return
	}

	var payloadDb models.Webhook
	payloadDb.ID = uuid.New()
	payloadDb.EventType = eventType
	payloadDb.SubscriptionID = subscriptionID
	payloadDb.CreatedAt = time.Now()
	payloadDb.UpdatedAt = time.Now()
	payloadDb.Payload = rawPayload
	payloadDb.Status = "pending"
	payloadDb.Delivered = false
	payloadDb.Retries = 0

	queueData := map[string]interface{}{
		"subscription_id": subscriptionID.String(),
		"payload":         json.RawMessage(rawPayload),
		"target_url":      subscription.TargetURL,
		"webhook_id":      payloadDb.ID,
	}

	queueDataBytes, err := json.Marshal(queueData)
	if err != nil {
		log.Println("Error marshaling queue data:", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to queue webhook",
			Data:    nil,
		})
		return
	}

	err = webhookHelper.InsertWebhook(&payloadDb)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to insert webhook into database: %v", err),
			Data:    nil,
		})
		return
	}

	err = queue.EnqueueWebhook(redis.RedisCtx, string(queueDataBytes))
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to queue webhook: %v", err),
			Data:    nil,
		})
		return
	}

	utils.WriteJson(w, http.StatusAccepted, utils.APIResponse{
		Status:  "success",
		Message: "Webhook accepted for processing",
		Data: map[string]interface{}{
			"webhook_id": payloadDb.ID,
		},
	})
}


func isValidEventType(events []string, eventType string) bool {
	for _, event := range events {
		if event == eventType {
			return true
		}
	}
	return false
}

// func verifySignatureRaw(rawPayload []byte, secret string, actualSignature string) bool {
// 	hash := hmac.New(sha256.New, []byte(secret))
// 	hash.Write(rawPayload)
// 	expectedSignature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
// 	return expectedSignature == actualSignature
// }


func GetAllWebHooksRequestOfSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionIDStr := mux.Vars(r)["subscriptionID"]
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "invalid subscription ID",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}

	webhooks, err := webhookHelper.GetWebhooksBySubscriptionID(subscriptionID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "failed to fetch webhooks for the subscription",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}
	

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "webhooks for the subscription fetched successfully",
		Data:    webhooks,
	})
}

func GetaWebHookbyID(w http.ResponseWriter, r *http.Request) {
	webhookIDStr := mux.Vars(r)["webhookID"]
	webhookID, err := uuid.Parse(webhookIDStr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "invalid webhook ID",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}

	webhook, err := webhookHelper.GetWebhookByID(webhookID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "failed to fetch webhook details",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "webhook fetched successfully",
		Data:    webhook,
	})
}




