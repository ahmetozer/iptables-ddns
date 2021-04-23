package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

//Lookup Resolve hosts
func (h Domain) Lookup() (string, error) {
	var nameserver string
	var err error
	if IsIPv4(h.Ns) {
		nameserver = h.Ns + ":53"
	} else if IsIPv6(h.Ns) {
		nameserver = "[" + h.Ns + "]:53"
	} else {
		return "", fmt.Errorf("nameserver is not valid")
	}
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, nameserver)
		},
	}
	if h.Qtype == "ip" {
		h.Qtype = "ip4"
	}
	ips, err := r.LookupIP(context.Background(), h.Qtype, h.Name)
	if err != nil {
		return "", err
	}

	if len(ips) != 1 {
		return "", fmt.Errorf("resolver expected one ip but get %v", len(ips))
	}
	return fmt.Sprintf("%v", ips[0]), nil
}
