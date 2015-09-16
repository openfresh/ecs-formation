package aws

import (
	"github.com/stormcat24/ecs-formation/logger"
	"strings"
	"time"
)

func isRateExceeded(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "Rate exceeded") {
		logger.Main.Errorf("AWS API Error: %s. Retry after 10 seconds.", err.Error())
		time.Sleep(15 * time.Second)
		return true
	}

	return false
}
