package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type errorArray []error

//Iptables add and delete mode for Iptables commands to related domain
func (obj Domain) Iptables(first_resolution string, second_resolution string) error {
	// Remove netfilter rules for old IPs
	var errA errorArray
	if first_resolution != "" {
		err := obj.IptablesExecuter("delete", first_resolution)
		if err != nil {
			errA = append(errA, err)
		}
	}
	if second_resolution != "" {
		err := obj.IptablesExecuter("add", second_resolution)
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

//IptablesExecuter
func (obj Domain) IptablesExecuter(mode string, address string) error {
	var iptablesCommand string
	if obj.Qtype == "ip" {
		iptablesCommand = "iptables"
	} else if obj.Qtype == "ip6" {
		iptablesCommand = "ip6tables"
	} else {
		panic(fmt.Sprintf("Unknown iptables mode %s %s", obj.Name, obj.Qtype))
	}
	file, err := os.Open(iptablesFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		// If its not commend and contains domain name
		if len(scanner.Text()) != 0 && (scanner.Text()[0] != 35 && DomainContains(scanner.Text(), obj.Name, obj.Qtype)) {
			if debug {
				fmt.Println("Line", lineNumber, " ", scanner.Text())
			}

			if mode == "add" {
				// Detect Add or inject
				rule1 := ""
				words := strings.Fields(scanner.Text())
				for _, arg := range words {
					if arg == "-I" {
						rule1 = "-I"
						break
					} else if arg == "-A" {
						rule1 = "-A"
						break
					}
				}

				if rule1 == "" {
					words = append([]string{"-A"}, words...)
				}

				// if debug {
				// 	fmt.Printf("%s\n", words)
				// }

				obj.replaceHostToAddr(&words, address)

				if debug {
					fmt.Printf("%s %v\n", iptablesCommand, words)
				}

				currentRuleHash, newRule := ruleHash(iptablesCommand, words)
				_, currentRuleExist := ddnsCurrentTable[currentRuleHash]
				if !currentRuleExist {
					err := iptables(iptablesCommand, words...)
					if err != nil {
						fmt.Printf("%s\n", err)
					} else {
						ddnsCurrentTable[currentRuleHash] = newRule
					}
				}
			} else if mode == "delete" {
				// Detect Add or inject
				words := strings.Fields(scanner.Text())

				// if debug {
				// 	fmt.Printf("%s\n", words)
				// }
				if words[0] == "-I" || words[0] == "-A" {
					words[0] = "-D"
				} else {
					words = append([]string{"-D"}, words...)
				}
				obj.replaceHostToAddr(&words, address)

				if debug {
					fmt.Printf("%s %v\n", iptablesCommand, words)
				}
				err := iptables(iptablesCommand, words...)
				if err == nil {
					// Remove rule from ddnsCurrentTable
					currentRuleHash, _ := ruleHash(iptablesCommand, words)
					_, currentRuleExist := ddnsCurrentTable[currentRuleHash]
					if currentRuleExist {
						delete(ddnsCurrentTable, currentRuleHash)
					}
				} else {
					fmt.Printf("%s\n", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("\tIptables function %s type %s mode %s address %s\n", obj.Name, obj.Qtype, mode, address)
	return nil
}

func (obj Domain) replaceHostToAddr(k *[]string, new string) {
	for i := range *k {
		if StartWithHost((*k)[i], obj.Name, obj.Qtype) {
			(*k)[i] = strings.Replace((*k)[i], obj.Name, new, 1)
		}
	}
}

func iptables(iptablesCommand string, words ...string) error {
	wait := []string{"-w", "100"}
	cmd := exec.Command(iptablesCommand, append(words, wait...)...)
	var cmdErr, cmdOut bytes.Buffer
	cmd.Stderr = &cmdErr
	cmd.Stdout = &cmdOut
	err := cmd.Run()
	if err != nil || cmdOut.String() != "" {
		return fmt.Errorf("%s %s: %s %s %s", iptablesCommand, words, err, cmdErr.String(), cmdOut.String())
	}
	return nil
}

func ruleHash(iptablesCommand string, obj []string) (string, []string) {
	temp_obj := make([]string, len(obj), cap(obj))
	copy(temp_obj, obj)
	temp_obj[0] = iptablesCommand
	h := sha1.New()
	h.Write([]byte(strings.Join(temp_obj, "")))
	return fmt.Sprintf("%x", h.Sum(nil)), temp_obj
}
