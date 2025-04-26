package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	logsHelper "github.com/harshgupta9473/webhookDelivery/helpers/logs"
	webhookHelper "github.com/harshgupta9473/webhookDelivery/helpers/webhook"
	"github.com/harshgupta9473/webhookDelivery/models"
	"github.com/harshgupta9473/webhookDelivery/workers/queue"
)

var MaxRetryAttempts = 3

var RetryIntervals = []time.Duration{
	10 * time.Second,
	30 * time.Second,
	60 * time.Second,
}

type WebhookTask struct {
	SubscriptionID string          `json:"subscription_id"`
	TargetURL      string          `json:"target_url"`
	Payload        json.RawMessage `json:"payload"`
	WebHookID      string          `json:"webhook_id"`
}

func DeliveryWorker() {
	log.Println("running in background")
	for {
		queueDataString, err := queue.DequeWebhookTask()
		log.Println("queue popped")
		if err != nil {
			log.Printf("Error popping task from queue: %v", err)
			continue
		}

		var task WebhookTask
		err = json.Unmarshal([]byte(queueDataString), &task)
		if err != nil {
			log.Printf("Failed to unmarshal queue data: %v", err)
			continue
		}

		subscriptionID, err := uuid.Parse(task.SubscriptionID)
		if err != nil {
			log.Printf("Invalid subscription ID: %v", err)
			continue
		}
		webHookID, err := uuid.Parse(task.WebHookID)
		if err != nil {
			log.Printf("Invalid webhook ID: %v", err)
			continue
		}

		ProcessWebHooks(task.Payload, subscriptionID, task.TargetURL, webHookID)
	}
}

func ProcessWebHooks(payload json.RawMessage, subscriptionID uuid.UUID, targetURL string, webHookID uuid.UUID) {

	err := deliverPayload(payload, targetURL, 1, webHookID, subscriptionID)
	if err == nil {
		log.Printf("Successfully delivered payload for subscription %s on first attempt", subscriptionID)
		err=webhookHelper.UpdateWebhookStatusAndDelivery(webHookID,"success",true,time.Now())
			if err!=nil{
				log.Println("error updating webhooks by id")
			}
		return
	}

	log.Printf("First attempt failed for subscription %s, starting retry in background: %v", subscriptionID, err)

	go func() {
		retryErr := deliverPayloadWithRetry(payload, targetURL, webHookID, subscriptionID)
		if retryErr != nil {
			log.Printf("Retry failed for subscription %s: %v", subscriptionID, retryErr)
			err=webhookHelper.UpdateWebhookStatusAndDelivery(webHookID,"failed",false,time.Now())
			if err!=nil{
				log.Println("error updating webhooks by id")
			}
		}else{
			err=webhookHelper.UpdateWebhookStatusAndDelivery(webHookID,"success",true,time.Now())
			if err!=nil{
				log.Println("error updating webhooks by id")
			}
		}
	}()
}

func deliverPayload(payload json.RawMessage, targetURL string, num int, webHookID uuid.UUID, subscriptionID uuid.UUID) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	
	logs := models.DeliveryLog{
		WebhookID:      webHookID,
		SubscriptionID: subscriptionID,
		TargetURL:      targetURL,
		AttemptNumber:  num,
		Timestamp:      time.Now(),
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		logs.ErrorDetails = fmt.Sprintf("Error delivering payload: %v", err)  
		log.Println(err)
		logs.Status = "failed"
		err = logsHelper.StoreDeliveryLog(logs)
		if err != nil {
			log.Printf("Failed to store delivery log: %v", err)
		}
		return fmt.Errorf("error delivering payload: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		logs.ErrorDetails = fmt.Sprintf("Error reading response body: %v", readErr)  
		log.Println("Error reading response body:", readErr)
	} else {
		fmt.Println("Response Body:", string(body))
	}

	logs.HTTPStatusCode = resp.StatusCode

	// Check for non-2xx response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logs.Status = "failed"
		logs.ErrorDetails = fmt.Sprintf("Non-2xx response: %d", resp.StatusCode)  // Capture error details here
		err = logsHelper.StoreDeliveryLog(logs)
		if err != nil {
			log.Printf("Failed to store delivery log: %v", err)
		}
		return fmt.Errorf("non-2xx response: %d", resp.StatusCode)
	}

	// On success, update status and store log
	logs.Status = "success"
	err = logsHelper.StoreDeliveryLog(logs)
	if err != nil {
		log.Printf("Failed to store delivery log: %v", err)
	}

	log.Printf("Successfully delivered webhook to %s", targetURL)
	return nil
}


func deliverPayloadWithRetry(payload json.RawMessage, targetURL string, webHookID uuid.UUID,subscriptionID uuid.UUID) error {
	var err error
	for attempt := 0; attempt < MaxRetryAttempts; attempt++ {
		log.Printf("Retry attempt %d for %s", attempt+1, targetURL)
		err = deliverPayload(payload, targetURL, attempt+1, webHookID,subscriptionID)
		if err == nil {
			log.Printf("Successfully delivered on retryattemt- %d (attempt %d)", attempt+1,attempt+2)
			return nil
		}
		log.Printf("Attempt %d failed: %v", attempt+2, err)
		// Only sleep if it's not the last retry attempt
		if attempt < MaxRetryAttempts-1 {
			wait := RetryIntervals[attempt]
			log.Printf("Waiting %v before retrying...", wait)
			time.Sleep(wait)
		}
	}

	return fmt.Errorf("failed to deliver payload after %d retries: %v", MaxRetryAttempts, err)
}
