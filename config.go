package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Defaults struct {
	Qtype    string `json:"type"`
	Ns       string `json:"ns"`
	Interval int    `json:"interval"`
}

type configStruct struct {
	Defaults `json:"defaults"`
	Domains  Domains `json:"domains"`
}

func LoadConfig() (configStruct, error) {
	defaultConf := configStruct{
		Defaults: Defaults{
			Qtype:    "ip",
			Ns:       "1.1.1.1",
			Interval: 300,
		},
	}
	configFileData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return defaultConf, fmt.Errorf(fmt.Sprintf("Err while reading the config.json file: %s", err))
	}

	err = json.Unmarshal([]byte(configFileData), &defaultConf)
	if err != nil {
		return defaultConf, fmt.Errorf(fmt.Sprintf("Err while parsing the config: %s", err))
	}

	if err = defaultConf.CheckDefaults(); err != nil {
		return defaultConf, err
	}

	return defaultConf, nil
}

func (obj configStruct) GetDefaults() []string {
	return []string{obj.Defaults.Qtype, obj.Defaults.Ns, fmt.Sprint(obj.Defaults.Interval)}
}

func (obj configStruct) PrintDefaults() string {
	defaults := obj.GetDefaults()
	return fmt.Sprintf("Defaults\ttype:%s \tns:%s\tinterval:%v\n", defaults[0], defaults[1], defaults[2])
}

func (obj configStruct) DuplicatedHostsCheck() []string {
	host_count := make(map[string]int)
	var duplicatedHosts []string
	for i := range obj.Domains {

		host_count[obj.Domains[i].Name]++
		if host_count[obj.Domains[i].Name] == 2 {
			duplicatedHosts = append(duplicatedHosts, obj.Domains[i].Name)
		}

	}
	return duplicatedHosts
}

func (obj configStruct) CheckDefaults() error {
	var errors []string

	if !(obj.Defaults.Qtype == "ip" || obj.Defaults.Qtype == "ip6") {
		errors = append(errors, "Type is not IP or IPv6: "+obj.Qtype)
	}

	if obj.Defaults.Interval < 300 {
		errors = append(errors, "Interval is too low: "+fmt.Sprint(obj.Interval))
	}

	if !(IsIPv4(obj.Defaults.Ns) || IsIPv6(obj.Defaults.Ns)) {
		errors = append(errors, "Name server is not valid IP address: "+obj.Ns)
	}
	if errors != nil {
		return fmt.Errorf(fmt.Sprintf("Defaults has a error: %s", errors))
	}
	return nil
}
