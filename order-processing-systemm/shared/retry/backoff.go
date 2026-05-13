package retry

import (
	"time"
)

func ExponentialBackoff(
	attempt int,
) time.Duration {

	baseDelay := time.Second

	return baseDelay * time.Duration(
		1<<attempt,
	)
}
