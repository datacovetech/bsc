package utils

import (
	"time"
)

func IntervalTask(itvl time.Duration, fn func()) chan struct{} {
	ticker := time.NewTicker(itvl)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fn()
			}
		}
	}()
	return done
}
