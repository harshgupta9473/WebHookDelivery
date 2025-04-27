package subscriptions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/webhookDelivery/cache"
	subscriptionDBHelper "github.com/harshgupta9473/webhookDelivery/helpers/subscriptions"
	"github.com/harshgupta9473/webhookDelivery/models"
	"github.com/harshgupta9473/webhookDelivery/utils"
)

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	sub.ID = uuid.New()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()
	if len(sub.EventTypes)==0{
		log.Println("Error creating subscription as there is no event type")
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Could not create subscription",
			Error:   "No event type provided",
		})
		return
	}

	err := subscriptionDBHelper.InsertSubscription(sub)
	
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Could not create subscription",
			Error:   err.Error(),
		})
		return
	}
	
	utils.WriteJson(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "Subscription created successfully",
		Data:    sub,
	})
}

func UpdateSubscription(w http.ResponseWriter, r *http.Request) {

	subscriptionIDStr := mux.Vars(r)["subscriptionID"]
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid subscription ID",
			Error:   err.Error(),
		})
		return
	}

	var updatedSub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&updatedSub); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	err = subscriptionDBHelper.UpdateSubscription(subscriptionID, &updatedSub)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to update subscription",
			Error:   err.Error(),
		})
		return
	}
	cache.DeleteCachedSubscriptionDetails(subscriptionID)

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Subscription updated successfully",
		Data: nil,
	})
}

func DeleteSubscription(w http.ResponseWriter, r *http.Request) {

	subscriptionIDStr := mux.Vars(r)["subscriptionID"]
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid subscription ID",
			Error:   err.Error(),
		})
		return
	}

	err = subscriptionDBHelper.DeleteSubscription(subscriptionID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to delete subscription",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Subscription deleted successfully",
	})
}

func GetSubscriptionDetails(w http.ResponseWriter, r *http.Request) {

	subscriptionIDStr := mux.Vars(r)["subscriptionID"]
	subscriptionID, err := uuid.Parse(subscriptionIDStr)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "Invalid subscription ID",
			Error:   err.Error(),
		})
		return
	}

	sub, err := subscriptionDBHelper.GetSubscriptionDetailsUsingDBorCache(subscriptionID)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Error retrieving subscription details",
			Error:   err.Error(),
		})
		return
	}

	if sub == nil {
		utils.WriteJson(w, http.StatusNotFound, utils.APIResponse{
			Status:  "error",
			Message: "Subscription not found",
			Error: "there no subscription related to the subscriptionID",
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Subscription retrieved successfully",
		Data:    sub,
	})
}

func GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	// Fetch all subscriptions using the helper function
	subscriptions, err :=subscriptionDBHelper.GetAllSubscriptions()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "failed to fetch subscriptions",
			Data:    nil,
			Error: err.Error(),
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "subscriptions fetched successfully",
		Data:    subscriptions,
	})
}
