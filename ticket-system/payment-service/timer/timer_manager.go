package timer

import (
	"log"
	"sync"
	"time"
)

var timers = make(map[int]chan bool)
var mu sync.Mutex

func Start(seatID int, duration time.Duration, onExpire func()) {
	mu.Lock()
	if _, exists := timers[seatID]; exists {
		log.Println("timer already running for seat:", seatID)
		mu.Unlock()
		return
	}

	stopChan := make(chan bool)
	timers[seatID] = stopChan

	mu.Unlock()

	log.Println("started hold timer for seat:", seatID)

	go func() {
		remaining := int(duration.Seconds())

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {

			case <-ticker.C:
				min := remaining / 60
				sec := remaining % 60

				log.Printf("Seat %d → %02d:%02d\n", seatID, min, sec)

				remaining--

				if remaining < 0 {
					log.Println("hold timer expired for seat:", seatID)
					onExpire()

					mu.Lock()
					delete(timers, seatID)
					mu.Unlock()

					return
				}

			case <-stopChan:
				log.Println("stopped hold timer for seat:", seatID)
				return
			}
		}
	}()
}

func Stop(seatID int) {
	mu.Lock()
	defer mu.Unlock()

	if ch, ok := timers[seatID]; ok {
		close(ch)
		delete(timers, seatID)
	}
}
