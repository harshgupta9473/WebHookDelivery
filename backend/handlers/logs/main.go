package logs

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	logsHelper "github.com/harshgupta9473/webhookDelivery/helpers/logs"
	"github.com/harshgupta9473/webhookDelivery/utils"
)

func GetLogsOfSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionIDStr := mux.Vars(r)["subscriptionID"]
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "invalid subscription ID",
			Data:   nil,
			Error: err.Error(),
		})
		return
	}

	logs, err := logsHelper.GetLogsBySubscriptionID(subscriptionID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "failed to fetch logs for the subscription",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "logs for the subscription fetched successfully",
		Data:    logs,
	})
}


func GetLogsOfWebhookID(w http.ResponseWriter, r *http.Request) {
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

	logs, err := logsHelper.GetLogsByWebhookID(webhookID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "failed to fetch logs for the webhook",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "logs for the webhook fetched successfully",
		Data:    logs,
	})
}

