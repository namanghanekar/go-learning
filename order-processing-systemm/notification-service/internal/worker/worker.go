package worker

import "fmt"

func Worker(id int, jobs <-chan NotificationJob) {

	for job := range jobs {
		fmt.Printf(
			"Worker %d processing notification for user %d\n",
			id,
			job.UserID,
		)
	}
}
