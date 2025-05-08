package domain

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type DomainNode struct {
	children      map[string]*DomainNode
	isEndOfDomain bool
}

type DomainTrie struct {
	root *DomainNode
}

func NewDomainTrie() *DomainTrie {
	return &DomainTrie{
		root: &DomainNode{
			children: make(map[string]*DomainNode),
		},
	}
}

func (t *DomainTrie) AddDomain(domain string) {
	parts := strings.Split(domain, ".")
	current := t.root

	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if _, exists := current.children[part]; !exists {
			current.children[part] = &DomainNode{
				children: make(map[string]*DomainNode),
			}
		}
		current = current.children[part]
	}
	current.isEndOfDomain = true
}

func (t *DomainTrie) MatchDomain(domain string) bool {
	parts := strings.Split(domain, ".")
	return t.matchDomainParts(parts, 0, t.root, false)
}

func (t *DomainTrie) matchDomainParts(parts []string, index int, node *DomainNode, foundComplete bool) bool {
	if index == len(parts) {
		return foundComplete
	}

	part := parts[len(parts)-1-index]

	if childNode, exists := node.children[part]; exists {
		newFoundComplete := foundComplete || childNode.isEndOfDomain
		if t.matchDomainParts(parts, index+1, childNode, newFoundComplete) {
			return true
		}
	}

	return foundComplete
}

type Manager struct {
	domainTrie *DomainTrie
}

func NewManager() *Manager {
	return &Manager{
		domainTrie: NewDomainTrie(),
	}
}

func (m *Manager) LoadDomains(domainFile string) error {
	file, err := os.Open(domainFile)
	if err != nil {
		return fmt.Errorf("failed to open domain file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			m.domainTrie.AddDomain(domain)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading domain file: %w", err)
	}

	return nil
}

func (m *Manager) IsChinaDomain(domain string) bool {
	return m.domainTrie.MatchDomain(domain)
}
