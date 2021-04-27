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
	debug            bool = false
	printConfig      bool = false
	keep             bool = false
	printVersion     bool = false
	printHelp        bool = false
	iptablesFile     string
	configFile       string
	GitTag           string = ""
	GitCommit        string = ""
	GitUrl           string = ""
	BuildTime        string = ""
	RunningEnv       string = ""
	ddnsCurrentTable map[string][]string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&keep, "keep", false, "Don't remove changes on exit")
	flag.BoolVar(&printVersion, "v", false, "Print version")

	if RunningEnv == "container" {
		flag.StringVar(&iptablesFile, "l", "/config/iptables.list", "Iptables rule list")
		flag.StringVar(&configFile, "f", "/config/config.json", "Program config file")
	} else {
		flag.StringVar(&iptablesFile, "l", "iptables.list", "Iptables rule list")
		flag.StringVar(&configFile, "f", "config.json", "Program config file")
	}
	flag.BoolVar(&printConfig, "p", false, "Prints configs per hosts")
	flag.BoolVar(&printHelp, "h", false, "Help")
	flag.Parse()

	if printVersion {
		pVersion()
	}
	if printHelp {
		flag.PrintDefaults()
		pVersion()
	}

	// Check the required programs
	errOnExit := 0
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
		errOnExit++
	}

	if _, err := os.Stat(iptablesFile); err != nil {
		fmt.Printf("Error while accesing %s:\t%s\n", iptablesFile, err)
		errOnExit++
	}

	if _, err := os.Stat(configFile); err != nil {
		fmt.Printf("Error while accesing %s:\t%s\n", configFile, err)
		errOnExit++
	}

	if !(os.Getenv("bDebug") == "true") && !checkCap("cap_net_admin") {
		fmt.Printf("Error program does not have cap_net_admin capabilities\n")
		if RunningEnv == "container" {
			fmt.Printf("execute container with \"--cap-add net_admin\" arg\n")
		}
		errOnExit++
	}

	if errOnExit > 0 {
		if errOnExit == 1 {
			fmt.Println("The program will be closed due to an error.")
		} else {
			fmt.Println("The program will be closed due to errors")
		}

		os.Exit(3)
	}

	ddnsCurrentTable = make(map[string][]string)
}
func main() {
	if debug {
		fmt.Printf("Debug: on\n")
	}

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
	if !keep {
		fmt.Println("Restoring Changes")
		for a := range ddnsCurrentTable {
			command := ddnsCurrentTable[a][0]
			ddnsCurrentTable[a][0] = "-D"
			if debug {
				fmt.Printf("%s %s %s\n", a, command, ddnsCurrentTable[a])
			}
			iptables(command, ddnsCurrentTable[a]...)
		}
	}
	fmt.Println("\nGood bye.")
}

func pVersion() {
	fmt.Printf("\n\tIptables-DDNS is a firewall management tool to keep update firewall rules with dynamic DNS.\n")

	if BuildTime != "" {
		fmt.Printf("\tProgram build date: %s\n", BuildTime)
	}

	if GitCommit != "" {
		fmt.Printf("\tCommmit: %s\n", GitCommit)
	}

	if GitTag != "" {
		fmt.Printf("\tTag: %s\n", GitTag)
	}

	if GitUrl != "" {
		fmt.Printf("\tRepo Url: %s\n", GitUrl)
	}

	os.Exit(0)
}
