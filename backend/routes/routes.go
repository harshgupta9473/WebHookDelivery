package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/webhookDelivery/handlers/logs"
	"github.com/harshgupta9473/webhookDelivery/handlers/subscriptions"
	"github.com/harshgupta9473/webhookDelivery/handlers/webhook"
)

func RegisterRoutes(r *mux.Router) {
	// subscription
	r.HandleFunc("/subscriptions",subscriptions.CreateSubscription).Methods("POST")
	r.HandleFunc("/subscriptions/{subscriptionID}",subscriptions.GetSubscriptionDetails).Methods("GET")
	r.HandleFunc("/subscriptions/{subscriptionID}",subscriptions.UpdateSubscription).Methods("PUT")
	r.HandleFunc("/subscriptions/{subscriptionID}",subscriptions.DeleteSubscription).Methods("DELETE")
	r.HandleFunc("/subscriptions",subscriptions.GetAllSubscriptions).Methods("GET")
	

	// ingestion into webhooks
	r.HandleFunc("/ingest/{subscriptionID}",webhook.IngestWebhook).Methods("POST")

	//logs
	r.HandleFunc("/logs/subscription/{subscriptionID}",logs.GetLogsOfSubscription).Methods("GET")
	r.HandleFunc("/logs/webhook/{webhookID}",logs.GetLogsOfWebhookID)

	//get webhooksrequests

	r.HandleFunc("/subscriptions/webhooks/{subscriptionID}",webhook.GetAllWebHooksRequestOfSubscription)
	r.HandleFunc("/requests/webhooks/{webhookID}",webhook.GetaWebHookbyID)
}