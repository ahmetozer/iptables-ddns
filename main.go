// +build linux

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	debug        bool = false
	iptablesFile string
	configFile   string
	printConfig  bool = false
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	if os.Getenv("runningenv") == "container" {
		flag.StringVar(&iptablesFile, "l", "/config/iptables.list", "Iptables rule list")
		flag.StringVar(&configFile, "-f", "/config/config.json", "Program config file")
	} else {
		flag.StringVar(&iptablesFile, "l", "iptables.list", "Iptables rule list")
		flag.StringVar(&configFile, "-f", "config.json", "Program config file")
	}
	flag.BoolVar(&printConfig, "-p", false, "Print config file")
	flag.Parse()

	// Check the required programs
	errOnExit := false
	requiredPrograms := []string{"iptables", "ip6tables"}
	var found int
	for i, s := range requiredPrograms {
		// Get path of required programs
		_, err := exec.LookPath(s)
		if err == nil {
			found++
		} else {
			fmt.Printf("Required program %v : %v cannot found.\n", i+1, s)
		}
	}
	if found != len(requiredPrograms) { //sh and df is must required. If is not found in software than exit.
		fmt.Printf("Please install required programs and re-execute this\n")
		errOnExit = true
	}

	if _, err := os.Stat(iptablesFile); err != nil {
		fmt.Printf("Error while accesing %s:\t%s\n", iptablesFile, err)
		errOnExit = true
	}

	if !checkCap("cap_net_admin") {
		fmt.Printf("Error program does not have cap_net_admin capabilities\n")
		if os.Getenv("runningenv") == "container" {
			fmt.Printf("execute container with \"--cap-add net_admin\" arg\n")
		}
		errOnExit = true
	}

	if errOnExit {
		os.Exit(3)
	}

}
func main() {
	programConf, err := LoadConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		fmt.Printf("config: loaded\n")
	}

	if printConfig || debug {
		fmt.Printf("%s", programConf.PrintDefaults())
	}

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
	} else {
		fmt.Printf("config: ok\n")
	}

	// Print configurations per domain
	if printConfig || debug {
		fmt.Printf("Configs for per hosts:\n")
		for i := range programConf.Domains {
			fmt.Printf("%s\n", programConf.Domains[i].String())
		}
		fmt.Println("")
	}

	// Start DDNS function
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
