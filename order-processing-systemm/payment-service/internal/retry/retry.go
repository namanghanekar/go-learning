package retry

import (
	"log"
	"time"
)

func Retry(
	attempts int,
	sleep time.Duration,
	fn func() error,
) error {

	var err error

	for i := 0; i < attempts; i++ {

		err = fn()

		if err == nil {
			return nil
		}

		log.Printf(
			"Retry %d failed: %v\n",
			i+1,
			err,
		)

		time.Sleep(sleep)
	}

	return err
}
