package retry

import (
	"fmt"
	"time"
)

func Retry(
	attempts int,
	fn func() error,
) error {

	var err error

	for i := 0; i < attempts; i++ {

		err = fn()

		if err == nil {
			return nil
		}

		delay := ExponentialBackoff(i)

		fmt.Println(
			"Retry attempt:",
			i+1,
			"waiting:",
			delay,
		)

		time.Sleep(delay)
	}

	return err
}
