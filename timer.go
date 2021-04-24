package main

import (
	"fmt"
	"time"
)

type schedulers []*time.Ticker

func (obj Domain) schedule() *time.Ticker {
	if debug {
		fmt.Printf("Schueduler starterd for %s:%s\n", obj.Name, time.Duration(obj.Interval)*time.Second)
	}

	//ticker := time.NewTicker(5 * time.Second)
	ticker := time.NewTicker(time.Duration(obj.Interval) * time.Second)
	go func() {
		first_resolution, second_resolution := "", ""
		for ; true; <-ticker.C {
			temp_resolution, err := obj.Lookup()
			if err == nil {
				second_resolution = temp_resolution
				if debug {
					fmt.Printf("%v %v ", obj.Name, temp_resolution)
				}
				if first_resolution != second_resolution {
					if debug {
						fmt.Printf("Addr changed")
					}
					first_resolution = second_resolution
					//Exec function
				} else if debug {
					fmt.Printf("Addr same")
				}
				if debug {
					fmt.Printf("\n")
				}
			} else {
				fmt.Printf("Error while querying host \"%s\": %s\n", obj.Name, err)
			}

		}
	}()
	return ticker
}
