package main

import (
	"fmt"
	"time"
)

type schedulers []*time.Ticker

func (obj Domain) schedule() *time.Ticker {
	fmt.Printf("Schueduler starterd for %s:%s\n", obj.Name, time.Duration(obj.Interval)*time.Second)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for ; true; <-ticker.C {
			// Exec function
		}
	}()
	return ticker
}
