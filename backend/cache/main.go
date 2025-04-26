package cache

import (
	"encoding/json"
	"log"
	"time"


	"github.com/google/uuid"
	"github.com/harshgupta9473/webhookDelivery/models"
	redisHelper "github.com/harshgupta9473/webhookDelivery/redis"
	redis "github.com/go-redis/redis/v8"
)

func CacheSubscriptionDetails(subscriptionID uuid.UUID, subscription models.Subscription) error {
	subscriptionData, err := json.Marshal(subscription)
	if err != nil {
		return err
	}

	err = redisHelper.RedisClient.Set(redisHelper.RedisCtx, subscriptionID.String(), subscriptionData, 1*time.Hour).Err()
	if err != nil {
		log.Printf("Error caching subscription details for %v: %v", subscriptionID, err)
		return err
	}
	return nil
}

func GetCachedSubscriptionDetails(subscriptionID uuid.UUID) (models.Subscription, error) {
	var subscription models.Subscription


	subscriptionData, err :=redisHelper.RedisClient.Get(redisHelper.RedisCtx, subscriptionID.String()).Result()
	if err == redis.Nil {
		log.Println("Not in Cache")
		return models.Subscription{}, err
	} else if err != nil {
		log.Printf("Error retrieving from Redis cache: %v", err)
		return models.Subscription{}, err
	} else {
		// Cache hit, deserialize the data
		err = json.Unmarshal([]byte(subscriptionData), &subscription)
		if err != nil {
			log.Printf("Error unmarshalling subscription data from cache: %v", err)
			return models.Subscription{}, err
		}
	}
	log.Println("found in cache")
	return subscription, nil
}
