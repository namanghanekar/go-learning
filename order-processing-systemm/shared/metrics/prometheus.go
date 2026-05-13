package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetricsServer() {

	http.Handle(
		"/metrics",
		promhttp.Handler(),
	)

	fmt.Println(
		"Metrics running on :2112",
	)

	http.ListenAndServe(
		":2112",
		nil,
	)
}
