# smolboydns
smolboydns is a lightweight DNS responder written in Go. It maps specific domains to IP addresses via a simple configuration file and logs incoming queries in real time.

## Quick Start
```bash
go build -o smolboydns smolboydns.go
echo "test.local. 10.10.10.10" > records.txt
sudo ./smolboydns records.txt
```
## Features
- Responds to A record queries from a static configuration file
- Logs each query with domain, resolved IP, and source IP
- Displays a live-updating table of requests in the console
- Single Go binary with no external services required

## Requirements
Go 1.20+
Root / administrator privileges to bind to UDP/53
Firewall rules allowing inbound UDP/53 traffic

## Installation
```bash
git clone <your-repo-url>
cd <repo>
go mod tidy
go build -o smolboydns smolboydns.go
```
## Usage
```bash
sudo ./smolboydns <config file>
```
### Config file format
The configuration file defines static DNS mappings. Each line contains a fully qualified domain name followed by an IPv4 address.
Rules:
        Domains must include a trailing dot (e.g. portal.example.com.)
        Domains should be lowercase
        Only A records are supported
Example records.txt:
```txt
example.com. 10.10.10.7
cdn.example.com.    10.10.10.8
owa.example.net.    172.16.10.5
```
## Deployment
Domains must be lowercase and end with a trailing dot
Only A records are supported
UDP only (TCP retries are not handled)
Config is loaded once at startup; restart to apply changes

