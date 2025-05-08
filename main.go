package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chareice/godnsproxy/dns"
)

func main() {
	var (
		domainFile     string
		port           int
		chinaServer    string
		trustDNSServer string
	)

	flag.StringVar(&domainFile, "f", "", "Path to domain file")
	flag.IntVar(&port, "p", 53, "DNS server port")
	flag.StringVar(&chinaServer, "c", "223.5.5.5", "China DNS server address")
	flag.StringVar(&trustDNSServer, "t", "https://1.1.1.1/dns-query", "Trust DNS server URL")
	flag.Parse()

	if domainFile == "" {
		fmt.Println("Error: domain file is required")
		flag.Usage()
		os.Exit(1)
	}

	server := dns.NewServer(&dns.Config{
		DomainFile:     domainFile,
		Port:           port,
		ChinaServer:    chinaServer,
		TrustDNSServer: trustDNSServer,
	})

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}
