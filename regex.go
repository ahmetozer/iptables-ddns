package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	/*
		Regex`s for checking input
	*/
	ipv6Regex   = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	ipv4Regex   = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
	domainRegex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z]$`
	portRegex   = "^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([1-9][0-9]{3})|([1-9][0-9]{2})|([1-9][0-9])|([1-9]))$"
)

//IsIPv4 is it valid IPv4 address
func IsIPv4(host string) bool {
	match, err := regexp.MatchString(ipv4Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

//IsIPv6 is it valid IPv6 address
func IsIPv6(host string) bool {
	match, err := regexp.MatchString(ipv6Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

//IsPort is it valid port number
func IsPort(port string) bool {
	match, err := regexp.MatchString(portRegex, port)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

//IsDomain check is valid host
func IsDomain(host string) bool {
	match, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func DomainContains(text string, search string, mode string) bool {
	words := strings.Fields(text)
	for _, arg := range words {
		if StartWithHost(arg, search, mode) {
			return true
		}
	}
	return false
}

func StartWithHost(text string, search string, mode string) bool {
	if mode == "ip" {
		match, err := regexp.MatchString("^"+search, text)
		if err != nil {
			fmt.Println(err)
			return false
		}
		return match
	} else {
		match, err := regexp.MatchString("^(\\[)?"+search, text)
		if err != nil {
			fmt.Println(err)
			return false
		}
		return match
	}
}
