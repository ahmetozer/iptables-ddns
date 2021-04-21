package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Defaults struct {
	Qtype    string `json:"type"`
	Ns       string `json:"ns"`
	Interval int    `json:"interval"`
}
type Domain struct {
	Name     string `json:"name"`
	Qtype    string `json:"type"`
	Ns       string `json:"ns"`
	Interval int    `json:"interval"`
}
type configStruct struct {
	Defaults `json:"defaults"`
	Domains  []Domain `json:"domains"`
}

func (obj *Domain) complete(Qtype string, Interval int, Ns string) {

	if obj.Qtype == "" {
		obj.Qtype = Qtype
	}
	if obj.Interval == 0 {
		obj.Interval = Interval
	}
	if obj.Ns == "" {
		obj.Ns = Ns
	}
}

func (k Domain) String() string {
	return fmt.Sprintf("%s\ttype:%s\tns:%s\tinterval:%v", k.Name, k.Qtype, k.Ns, k.Interval)
}
func main() {
	programConf := configStruct{
		Defaults: Defaults{
			Qtype:    "ip",
			Ns:       "1.1.1.1",
			Interval: 300,
		},
	}

	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("Err while reading the config.json file: %s", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(dat), &programConf)
	if err != nil {
		fmt.Printf("Err while parsing the config: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Default:\ttype:%s\tns:%s\tinterval:%v\n", programConf.Defaults.Qtype, programConf.Defaults.Ns, programConf.Defaults.Interval)

	host_count := make(map[string]int)
	var duplicatedHosts []string
	for i, _ := range programConf.Domains {

		host_count[programConf.Domains[i].Name]++
		if host_count[programConf.Domains[i].Name] == 2 {
			duplicatedHosts = append(duplicatedHosts, programConf.Domains[i].Name)
		}

	}
	if len(duplicatedHosts) != 0 {
		fmt.Printf("ERR: Some hosts are duplicated\n%s\n", duplicatedHosts)
		os.Exit(1)
	}
	fmt.Printf("Configs for per hosts:\n")
	for i, _ := range programConf.Domains {
		programConf.Domains[i].complete(programConf.Defaults.Qtype, programConf.Defaults.Interval, programConf.Defaults.Ns)
	}

	for i, _ := range programConf.Domains {
		fmt.Printf("%s\n", programConf.Domains[i].String())
	}
}
