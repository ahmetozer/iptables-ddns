package main

import "fmt"

type errorArray []error

//Iptables add and delte mode for Iptables commands to related domain
func (obj Domain) Iptables(first_resolution string, second_resolution string) error {
	// Remove netfilter rules for old IPs
	var errA errorArray
	if first_resolution != "" {
		err := obj.IptablesExecute("delete", first_resolution)
		if err != nil {
			errA = append(errA, err)
		}
	}
	if second_resolution != "" {
		err := obj.IptablesExecute("add", second_resolution)
		if err != nil {
			errA = append(errA, err)
		}
	} else {
		return fmt.Errorf("the second_resolution is not to be empty")
	}

	if errA != nil {
		return fmt.Errorf("%s", errA)
	}
	return nil
}

//IptablesExecute
func (obj Domain) IptablesExecute(mode string, address string) error {
	return nil
}
