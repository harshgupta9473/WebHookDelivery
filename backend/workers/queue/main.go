package queue

import (
	"context"
	"fmt"

	"github.com/harshgupta9473/webhookDelivery/redis"
)

func EnqueueWebhook(ctx context.Context, payload string) error {

	err := redis.RedisClient.LPush(redis.RedisCtx, "webhookQueue", payload).Err()
	if err != nil {
		return fmt.Errorf("failed to queue task: %v", err)
	}
	return nil
}

func DequeWebhookTask() (string, error) {
	// Use BRPOP to block and pop a task from the "webhookQueue" list
	result, err := redis.RedisClient.BRPop(redis.RedisCtx, 0, "webhookQueue").Result()
	if err != nil {
		return "", fmt.Errorf("failed to process task: %v", err)
	}
	return result[1], nil
}
