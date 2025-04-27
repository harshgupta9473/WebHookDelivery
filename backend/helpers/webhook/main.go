package webhookHelper

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/harshgupta9473/webhookDelivery/db"
	"github.com/harshgupta9473/webhookDelivery/models"
)

func InsertWebhook(webhook *models.Webhook) error {

	query := `INSERT INTO webhooks (id, event_type, payload, subscription_id, status, delivered, retries, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// Execute the query
	_, err := db.DB.Exec(query, webhook.ID, webhook.EventType, webhook.Payload, webhook.SubscriptionID, webhook.Status, webhook.Delivered, webhook.Retries, webhook.CreatedAt, webhook.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting webhook: %v", err)
		return err
	}

	log.Printf("Webhook inserted successfully: %v", webhook.ID)
	return nil
}

func UpdateWebhookStatusAndDelivery(webhookID uuid.UUID, status string, delivered bool, updatedAt time.Time,num int) error {
	query := `UPDATE webhooks 
	          SET status = $1, delivered = $2, updated_at = $3 ,retries=$4,
	          WHERE id = $5`

	_, err := db.DB.Exec(query, status, delivered, updatedAt,num, webhookID)
	if err != nil {
		log.Printf("Error updating webhook (ID: %v): %v", webhookID, err)
		return err
	}

	log.Printf("Webhook updated successfully (ID: %v)", webhookID)
	return nil
}

func GetWebhookByID(webhookID uuid.UUID) (*models.Webhook, error) {

	query := `SELECT id, event_type, payload, subscription_id, status, delivered, created_at, updated_at, retries 
	          FROM webhooks WHERE id = $1`

	row := db.DB.QueryRow(query, webhookID)

	var webhook models.Webhook

	err := row.Scan(&webhook.ID, &webhook.EventType, &webhook.Payload, &webhook.SubscriptionID, &webhook.Status, &webhook.Delivered, &webhook.CreatedAt, &webhook.UpdatedAt, &webhook.Retries)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no webhook found with ID %v", webhookID)
		}
		log.Printf("Error retrieving webhook: %v", err)
		return nil, err
	}

	return &webhook, nil
}

func GetWebhooksBySubscriptionID(subscriptionID uuid.UUID) ([]*models.Webhook, error) {
	query := `SELECT id, event_type, payload, subscription_id, status, delivered, created_at, updated_at, retries 
	          FROM webhooks WHERE subscription_id = $1`

	rows, err := db.DB.Query(query, subscriptionID)
	if err != nil {
		log.Printf("Error querying webhooks by subscription ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var webhooks []*models.Webhook

	for rows.Next() {
		var webhook models.Webhook
		err := rows.Scan(
			&webhook.ID,
			&webhook.EventType,
			&webhook.Payload,
			&webhook.SubscriptionID,
			&webhook.Status,
			&webhook.Delivered,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
			&webhook.Retries,
		)
		if err != nil {
			log.Printf("Error scanning webhook row: %v", err)
			return nil, err
		}
		webhooks = append(webhooks, &webhook)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return webhooks, nil
}
