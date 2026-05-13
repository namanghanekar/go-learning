package utils

import (
	"fmt"
	"time"
)

func GenerateTrackingID(
	orderID int,
) string {

	return fmt.Sprintf(
		"TRK-%d-%d",
		orderID,
		time.Now().Unix(),
	)
}
