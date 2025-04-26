package logsCleanup

import (
	"log"
	"time"

	logsHelper "github.com/harshgupta9473/webhookDelivery/helpers/logs"
)

func RunCleanupTask() {

	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("starting cleanup task")
		logsHelper.CleanupOldLogs()
	}
}
