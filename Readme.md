# iptables-ddns

Iptables-DDNS is a firewall management tool to keep update firewall rules with dynamic DNS.

## How it works ?

The application needs two configuration files which are config.json and iptables.list.
Config list includes default configurations for the application such as name server, default interval and default mode (ip,ip6) and hosts.  
After loading the config file at startup, the system starts to update functions per domain. Each update function resolves domains and if there is a change, functions inspect the iptables.list and if it`s related to domain, the function executes iptables to apply new rules and removes old rules.

## Why ?

Many devices such as modem and routers have DDNS support even if too old. Many clients don`t have a static IP and their addresses change in every new PPPoE session. With this, you can allow requests to services in bare port or IP without using proxy or VPN tools.

Another development reason for this system is to expose different services on the same port per IP addresses with DNAT rule.

## Configuration

In order for the system to work, you need to make settings on two different files.

### config.json

#### Defaults

Your domain addresses and their configurations are located in this file.  
The program has a built-in configuration but you can change these defaults with the `defaults` block in config.json  

Built-in defaults

```json
"defaults": {
       "type": "ip",
       "interval": 300,
       "ns": "1.1.1.1"
}
```

#### Domain

You can change the configuration per domain in the `domains` block.
You have to create an object per domain in the `domains` array to work and the domain name must be indicated at the name variable.  
**Note:** Non indicated domains in iptables.list will not work.  
You can change the type of domain with the `type` variable to “ip” or “ip6”. Also, the public recursive name server caches the request. If you want to update your domains much faster, you can directly give your domain's name server in `ns` variable and the system is not affected by the cache.

```json
{
           "name": "ddns.example.domain.test",
           "type": "ip",
           "ns":"203.0.113.56"
}
```

#### Example config.json

```json
{
   "defaults": {
       "type": "ip",
       "interval": 300
   },
   "domains": [
       {
           "name": "ddns.example.domain.test",
           "type": "ip",
           "ns":"203.0.113.56"
       },
        {
           "name": "2ddns.example.domain.test"
       },
       {
           "name": "3ddns.example.domain.test",
           "ns": "192.0.2.57"
       },
       {
           "name": "4ddns.example.domain.test",
           "ns": "192.0.2.57",
           "type":"ip6"
       }
   ]
}
```

### iptables.list

Your iptables rules are located in this file. You can do a block, allow or other iptables related actions within this file.  
Your configuration must be like an `iptables -S` command output.

**Note:** If your rule is in the table, don’t forget to add table argument `-t` with table name like `-t nat` in your rule line.

```txt
# Example Iptables List

# Route incoming requests to container
-A PREROUTING -i eth0 -s ddns.example.domain.test -p tcp  --dport 443 -j DNAT --to-destination 172.23.0.50:8443 -t nat
-A PREROUTING -i eth0 -s 2ddns.example.domain.test -p tcp  --dport 443 -j DNAT --to-destination 172.23.0.52:8443 -t nat

#  Forward ip to remote address
-A PREROUTING -i lan1 -d 172.23.0.50 -j DNAT --to-destination ddns.example.domain.test -t nat

# Allow SSH
-A INPUT -p tcp -s 3ddns.example.domain.test --dport 22 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
-A OUTPUT -p tcp -d 3ddns.example.domain.test  --sport 22 -m conntrack --ctstate ESTABLISHED -j ACCEPT

# Allow IP
-A INPUT -i eth0 -s 4ddns.example.domain.test -j ACCEPT
```

## Run

You can run this software inside the container on the host network or build it from the source and execute it on your system.

### Arguments

```md
Usage of iptables-ddns:
  -debug
        Debug mode
  -f string
        Program config file (default "/config/config.json")
  -keep
        Don`t remove changes on exit
  -l string
        Iptables rule list (default "/config/iptables.list")
  -p    Prints configs per hosts
  -v    Print version
  ```

### Container

After configuration in iptables.list and config.json, start the container with net_admin capabilities and mount your configuration folder to the container.

```bash
docker run -it -d --name iptables-ddns --rm --cap-add=net_admin --network host  -v /data/config/iptables-ddns/:/config/ ghcr.io/ahmetozer/iptables-ddns:latest
```

### From source

This system requires iptables and ip6tables. If you don't have it on the system, please install it before the execution.

```bash
go get -v github.com/ahmetozer/iptables-ddns
iptables-ddns
```
