package subscriptionDBHelper

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/harshgupta9473/webhookDelivery/cache"
	"github.com/harshgupta9473/webhookDelivery/db"
	"github.com/harshgupta9473/webhookDelivery/models"
	"github.com/lib/pq"
)

func InsertSubscription(subscription models.Subscription) error {

	query := `
        INSERT INTO subscriptions (id, target_url, secret, event_types)
        VALUES ($1, $2, $3, $4)
    `

	_, err := db.DB.Exec(query, subscription.ID, subscription.TargetURL, subscription.Secret,
		pq.Array(subscription.EventTypes))
	if err != nil {
		log.Printf("Error inserting subscription: %v", err)
		return err
	}
	return nil
}

func GetSubscriptionDetails(subscriptionID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription

	query := `SELECT id, target_url, secret, event_types, created_at, updated_at FROM subscriptions WHERE id = $1`
	row := db.DB.QueryRow(query, subscriptionID)

	var eventTypes pq.StringArray
	err := row.Scan(&subscription.ID, &subscription.TargetURL, &subscription.Secret, &eventTypes, &subscription.CreatedAt, &subscription.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No subscription found with given subscriptionID")
		}
		log.Printf("Error fetching subscription details: %v", err)
		return nil, err
	}

	subscription.EventTypes = []string(eventTypes)

	return &subscription, nil
}

func UpdateSubscription(subscriptionID uuid.UUID, data *models.Subscription) error {
	eventTypes := pq.StringArray(data.EventTypes)

	query := `
        UPDATE subscriptions 
        SET target_url = $1, secret = $2, event_types = $3, updated_at = NOW()
        WHERE id = $4
    `

	result, err := db.DB.Exec(query, data.TargetURL, data.Secret, eventTypes, subscriptionID)
	if err != nil {
		log.Printf("Error updating subscription: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {

		log.Printf("No subscription found with ID: %s", subscriptionID)
		return fmt.Errorf("No subscription found with ID: %s", subscriptionID)
	}

	return nil
}

func DeleteSubscription(subscriptionID uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := db.DB.Exec(query, subscriptionID)
	if err != nil {
		log.Printf("Error deleting subscription: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {

		log.Printf("No subscription found with ID: %s", subscriptionID)
		return fmt.Errorf("No subscription found with ID: %s", subscriptionID)
	}

	return nil
}

func GetAllSubscriptions() ([]*models.Subscription, error) {
	query := `SELECT id, target_url, secret, event_types, created_at, updated_at FROM subscriptions`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying subscriptions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription

	for rows.Next() {
		var subscription models.Subscription
		var eventTypes pq.StringArray
		err := rows.Scan(
			&subscription.ID,
			&subscription.TargetURL,
			&subscription.Secret,
			&eventTypes,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning subscription details: %v", err)
			return nil, err
		}
		subscription.EventTypes = []string(eventTypes)
		subscriptions = append(subscriptions, &subscription)
	}

	// Check if iteration produced any error
	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	// 👇 Extra check: no subscriptions found
	if len(subscriptions) == 0 {
		return nil, fmt.Errorf("no subscriptions found")
	}

	return subscriptions, nil
}

func GetSubscriptionDetailsUsingDBorCache(subscriptionID uuid.UUID) (*models.Subscription, error) {
	subscriptionFromCache, err := cache.GetCachedSubscriptionDetails(subscriptionID)
	if err == nil {
		return &subscriptionFromCache,nil
	}

	subscriptionFromDB,err:=GetSubscriptionDetails(subscriptionID)
	if err!=nil{
		return &models.Subscription{},err
	}
	err=cache.CacheSubscriptionDetails(subscriptionID,*subscriptionFromDB)
	if err!=nil{
		log.Println("error caching the subscription details")
	}
	return subscriptionFromDB ,nil
}
