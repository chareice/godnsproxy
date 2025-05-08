package dns

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/chareice/godnsproxy/domain"
)

type Config struct {
	DomainFile     string
	Port           int
	ChinaServer    string
	TrustDNSServer string
}

type Server struct {
	config     *Config
	domainMgr  *domain.Manager
	udpServer  *net.UDPConn
	httpClient *http.Client
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	// Initialize domain manager
	s.domainMgr = domain.NewManager()
	if err := s.domainMgr.LoadDomains(s.config.DomainFile); err != nil {
		return fmt.Errorf("failed to load domains: %w", err)
	}

	// Start UDP server
	addr := &net.UDPAddr{Port: s.config.Port}
	var err error
	s.udpServer, err = net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP port %d: %w", s.config.Port, err)
	}
	defer s.udpServer.Close()

	fmt.Printf("DNS server listening on port %d\n", s.config.Port)

	buf := make([]byte, 512)
	for {
		n, addr, err := s.udpServer.ReadFromUDP(buf)
		if err != nil {
			return fmt.Errorf("error reading from UDP: %w", err)
		}

		go s.handleRequest(buf[:n], addr)
	}
}

func (s *Server) handleRequest(query []byte, addr *net.UDPAddr) {
	start := time.Now()
	domain := getDomainFromDnsQuery(query)

	log.Printf("[%s] RECV %s query %s",
		time.Now().Format("15:04:05"),
		addr.IP.String(),
		domain)

	isChinaDomain := s.domainMgr.IsChinaDomain(domain)

	var resp []byte
	var err error
	var dnsType string

	if isChinaDomain {
		dnsType = "CHINA"
	} else {
		dnsType = "DOH"
	}

	log.Printf("[%s] PROC %s using %s DNS",
		time.Now().Format("15:04:05"),
		domain,
		dnsType)

	if isChinaDomain {
		resp, err = s.forwardToUDP(query)
	} else {
		resp, err = s.forwardToDoH(query)
	}

	duration := time.Since(start)
	status := "success"
	if err != nil {
		status = "failed"
		log.Printf("[%s] DONE %s %s %v",
			time.Now().Format("15:04:05"),
			domain,
			status,
			duration)
		fmt.Printf("Error handling request: %v\n", err)
		return
	}

	log.Printf("[%s] DONE %s %s %v",
		time.Now().Format("15:04:05"),
		domain,
		status,
		duration)

	if _, err := s.udpServer.WriteToUDP(resp, addr); err != nil {
		log.Printf("[%s] %s %s failed_to_send %v",
			time.Now().Format("15:04:05"),
			domain,
			dnsType,
			duration)
		fmt.Printf("Error sending response: %v\n", err)
	}
}

func (s *Server) forwardToUDP(query []byte) ([]byte, error) {
	conn, err := net.Dial("udp", s.config.ChinaServer+":53")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to China DNS: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(query); err != nil {
		return nil, fmt.Errorf("failed to send query to China DNS: %w", err)
	}

	resp := make([]byte, 512)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from China DNS: %w", err)
	}

	return resp[:n], nil
}

func (s *Server) forwardToDoH(query []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", s.config.TrustDNSServer, bytes.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create DoH request: %w", err)
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send DoH request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoH server returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func getDomainFromDnsQuery(msg []byte) string {
	if len(msg) < 13 {
		return "[malformed_dns_query]"
	}

	questionSection := msg[12 : len(msg)-4]
	var labels []string
	offset := 0

	for offset < len(questionSection) {
		length := int(questionSection[offset])
		offset++

		if length == 0 {
			break
		}

		// Check for DNS pointer (compression)
		if (length & 0xc0) == 0xc0 {
			// Skip pointer for now
			break
		}

		if offset+length > len(questionSection) {
			return "[malformed_domain]"
		}

		label := string(questionSection[offset : offset+length])
		labels = append(labels, label)
		offset += length
	}

	if len(labels) == 0 {
		// Fallback to simple parsing
		var fallback string
		for i := 0; i < len(questionSection) && questionSection[i] != 0; i++ {
			fallback += string(questionSection[i])
		}
		if fallback == "" {
			return "[unknown_domain]"
		}
		return fallback
	}

	return strings.Join(labels, ".")
}
