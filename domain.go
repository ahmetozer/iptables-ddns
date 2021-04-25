package main

import (
	"fmt"
)

type Domain struct {
	Name     string `json:"name"`
	Qtype    string `json:"type"`
	Ns       string `json:"ns"`
	Interval int    `json:"interval"`
}
type Domains []Domain

func (obj *Domain) Complete(Qtype string, Interval int, Ns string) {

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

//check Control domain properties
func (obj *Domain) Check() []string {
	var errors []string

	if !IsDomain(obj.Name) {
		errors = append(errors, "Domain is not valid: "+obj.Name)
	}

	if !(obj.Qtype == "ip" || obj.Qtype == "ip6") {
		errors = append(errors, "Type is not IP or ip6: "+obj.Qtype)
	}

	if obj.Interval < 300 {
		errors = append(errors, "Interval is too low: "+fmt.Sprint(obj.Interval))
	}

	if !(IsIPv4(obj.Ns) || IsIPv6(obj.Ns)) {
		errors = append(errors, "Name server is not valid IP address: "+obj.Ns)
	}

	return errors
}

//String Domain to string
func (k Domain) String() string {
	return fmt.Sprintf("%s type:%s ns:%s interval:%v", k.Name, k.Qtype, k.Ns, k.Interval)
}

func (k Domains) Check() [][]string {
	var stat [][]string
	for i := range k {
		if control := k[i].Check(); control != nil {
			stat = append(stat, append([]string{k[i].Name}, control...))
		}
	}
	return stat
}
