package logsHelper

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/harshgupta9473/webhookDelivery/db"
	"github.com/harshgupta9473/webhookDelivery/models"
)

func StoreDeliveryLog(log models.DeliveryLog) error {

	query := `INSERT INTO delivery_logs (webhook_id, subscription_id, target_url, attempt_number, status, http_status_code, error_details, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.DB.Exec(query, log.WebhookID, log.SubscriptionID, log.TargetURL, log.AttemptNumber, log.Status, log.HTTPStatusCode, log.ErrorDetails, log.Timestamp)

	if err != nil {
		return fmt.Errorf("error inserting delivery log: %v", err)
	}

	return nil
}

const RetentionPeriod = 72 * time.Hour //  logs for 72 hours

func CleanupOldLogs() {
	// Delete logs older than the retention period
	_, err := db.DB.Exec(`
		DELETE FROM delivery_logs WHERE timestamp < NOW() - INTERVAL '72 HOURS'`)
	if err != nil {
		log.Printf("Error cleaning up old logs: %v", err)
		return
	}

	log.Println("Successfully cleaned up old delivery logs")
}

func GetLogsByWebhookID(webhookID uuid.UUID) ([]models.DeliveryLog, error) {
	query := `SELECT id, webhook_id, subscription_id, target_url, attempt_number, status, http_status_code, error_details, timestamp
	          FROM delivery_logs WHERE webhook_id = $1 ORDER BY timestamp DESC`

	rows, err := db.DB.Query(query, webhookID)
	if err != nil {
		return nil, fmt.Errorf("error fetching delivery logs: %v", err)
	}
	defer rows.Close()

	var logs []models.DeliveryLog

	for rows.Next() {
		var log models.DeliveryLog
		if err := rows.Scan(
			&log.ID,
			&log.WebhookID,
			&log.SubscriptionID,
			&log.TargetURL,
			&log.AttemptNumber,
			&log.Status,
			&log.HTTPStatusCode,
			&log.ErrorDetails,
			&log.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("error scanning delivery log row: %v", err)
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs found for webhook ID %v", webhookID)
	}
	return logs, nil
}

func GetLogsBySubscriptionID(subscriptionID uuid.UUID) ([]models.DeliveryLog, error) {
	query := `SELECT id, webhook_id, subscription_id, target_url, attempt_number, status, http_status_code, error_details, timestamp
	          FROM delivery_logs WHERE subscription_id = $1 ORDER BY timestamp DESC`

	rows, err := db.DB.Query(query, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("error fetching delivery logs: %v", err)
	}
	defer rows.Close()

	var logs []models.DeliveryLog

	for rows.Next() {
		var log models.DeliveryLog
		if err := rows.Scan(
			&log.ID,
			&log.WebhookID,
			&log.SubscriptionID,
			&log.TargetURL,
			&log.AttemptNumber,
			&log.Status,
			&log.HTTPStatusCode,
			&log.ErrorDetails,
			&log.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("error scanning delivery log row: %v", err)
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}
	if len(logs) == 0 {
		return nil, fmt.Errorf("no logs found for subscription ID %v", subscriptionID)
	}

	return logs, nil
}
