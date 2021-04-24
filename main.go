// +build linux

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	debug bool = false
)

func init() {
	flag.BoolVar(&debug, "D", false, "Debug mode")
	flag.Parse()
}
func main() {
	programConf, err := LoadConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	fmt.Printf("%s", programConf.PrintDefaults())

	// Duplicated hosts check
	if duplicatedHosts := programConf.DuplicatedHostsCheck(); duplicatedHosts != nil {
		fmt.Printf("ERR: Some hosts are duplicated\n%s\n", duplicatedHosts)
		os.Exit(1)
	}

	// Complete configurations
	for i := range programConf.Domains {
		programConf.Domains[i].Complete(programConf.Defaults.Qtype, programConf.Defaults.Interval, programConf.Defaults.Ns)
	}

	// Check configuration per domain
	if stat := programConf.Domains.Check(); stat != nil {
		for i := range stat {
			fmt.Printf("Config error on domain \"%s\"", stat[i][0])
			for k := 1; k < len(stat[i]); k++ {
				fmt.Printf("\n\t %s", stat[i][k])
			}
			fmt.Print("\n")
		}
		os.Exit(1)
	}
	// Print configurations per domain
	fmt.Printf("Configs for per hosts:\n")
	for i := range programConf.Domains {
		fmt.Printf("%s\n", programConf.Domains[i].String())
	}

	// Start DDNS function
	fmt.Println("")
	var schedules schedulers
	for i := range programConf.Domains {
		schedules = append(schedules, programConf.Domains[i].schedule())
		defer schedules[i].Stop()
		// To avoid many DNS queries and nftables changes hits to rate limit
		time.Sleep(500 * time.Millisecond)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Schedulers are started.")
	<-sigs
	fmt.Println("\nGood bye.")
}
