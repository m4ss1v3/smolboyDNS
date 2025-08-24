package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var records = map[string]string{}

// Channel for sending DNS request logs
var requestChannel = make(chan string, 100)

func parseConfig(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			records[fields[0]] = fields[1]
		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range m.Question {
			switch q.Qtype {
			case dns.TypeA:
				ip := records[q.Name]
				if ip != "" {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
					// Extract the source IP address
					sourceIP := w.RemoteAddr().(*net.UDPAddr).IP.String()
					// Extract the requested domain and IP address
					requestedDomain := q.Name
					requestedIP := ip
					
					// Send request info to the channel for logging
					requestChannel <- fmt.Sprintf("%s -> %s (%s)", requestedDomain, requestedIP, sourceIP)
				}
			}
		}
	}

	w.WriteMsg(m)
}

func printRequests() {
	// To store records temporarily for printing the live feed
	recordsTable := make(map[string]string)

	for {
		select {
		case request := <-requestChannel:
			// On receiving a new request, update the table
			// You can add more sophisticated logic here if you want to keep track of old records
			parts := strings.Split(request, " -> ")
			if len(parts) == 3 {
				domain := parts[0]
				ip := parts[1]
				sourceIP := parts[2]

				// Update the records map with domain, IP, and source IP
				recordsTable[domain] = fmt.Sprintf("%s (%s)", ip, sourceIP)
			}
		}

		// Clear the terminal and print the updated table
		fmt.Print("\033[H\033[2J") // ANSI escape sequence to clear the screen
		fmt.Println("+------------+-------------+-----------------+-----------------+")
		fmt.Println("| Domain     | IP Address  | Source IP       | Request         |")
		fmt.Println("+------------+-------------+-----------------+-----------------+")

		for domain, ip := range recordsTable {
			fields := strings.Split(ip, " (")
			if len(fields) == 2 {
				requestedIP := fields[0]
				sourceIP := fields[1][:len(fields[1])-1]
				fmt.Printf("| %-10s | %-12s | %-15s | %-15s |\n", domain, requestedIP, sourceIP, ip)
			}
		}

		fmt.Println("+------------+-------------+-----------------+-----------------+")
		time.Sleep(1 * time.Second) // Update every second
	}
}

func printBanner() {
	// Your existing banner code
	fmt.Println("+---------------------------------------------------------------------------------------+")
	fmt.Println("|  ██████  ███▄ ▄███▓ ▒█████   ██▓     ▄▄▄▄    ▒█████   ██▓▓█████▄  ███▄    █   ██████  |")
	fmt.Println("|▒██    ▒ ▓██▒▀█▀ ██▒▒██▒  ██▒▓██▒    ▓█████▄ ▒██▒  ██▒▓██▒▒██▀ ██▌ ██ ▀█   █ ▒██    ▒  |")
	fmt.Println("|░ ▓██▄   ▓██    ▓██░▒██░  ██▒▒██░    ▒██▒ ▄██▒██░  ██▒▒██▒░██   █▌▓██  ▀█ ██▒░ ▓██▄    |")
	fmt.Println("|  ▒   ██▒▒██    ▒██ ▒██   ██░▒██░    ▒██░█▀  ▒██   ██░░██░░▓█▄   ▌▓██▒  ▐▌██▒  ▒   ██▒ |")
	fmt.Println("|▒██████▒▒▒██▒   ░██▒░ ████▓▒░░██████▒░▓█  ▀█▓░ ████▓▒░░██░░▒████▓ ▒██░   ▓██░▒██████▒▒ |")
	fmt.Println("|▒ ▒▓▒ ▒ ░░ ▒░   ░  ░░ ▒░▒░▒░ ░ ▒░▓  ░░▒▓███▀▒░ ▒░▒░▒░ ░▓   ▒▒▓  ▒ ░ ▒░   ▒ ▒ ▒ ▒▓▒ ▒ ░ |")
	fmt.Println("|░ ░▒  ░ ░░  ░      ░  ░ ▒ ▒░ ░ ░ ▒  ░▒░▒   ░   ░ ▒ ▒░  ▒ ░ ░ ▒  ▒ ░ ░░   ░ ▒░░ ░▒  ░ ░ |")
	fmt.Println("|░  ░  ░  ░      ░   ░ ░ ░ ▒    ░ ░    ░    ░ ░ ░ ░ ▒   ▒ ░ ░ ░  ░    ░   ░ ░ ░  ░  ░   |")
	fmt.Println("|      ░         ░       ░ ░      ░  ░ ░          ░ ░   ░     ░             ░       ░   |")
	fmt.Println("|                                           ░               ░                           |")
	fmt.Println("| Developed by M$ss1v3                                                                  |")
	fmt.Println("+---------------------------------------------------------------------------------------+")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config file>\n", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	parseConfig(configFile)

	printBanner()

	// Start a goroutine for printing the request feed
	go printRequests()

	// Start the DNS server
	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{Addr: ":53", Net: "udp"}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}
