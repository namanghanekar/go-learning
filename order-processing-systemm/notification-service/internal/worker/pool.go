package worker

func StartWorkerPool(workerCount int) chan NotificationJob {

	jobs := make(chan NotificationJob, 100)

	for i := 1; i <= workerCount; i++ {
		go Worker(i, jobs)
	}

	return jobs
}
