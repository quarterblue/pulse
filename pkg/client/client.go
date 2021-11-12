package client

import (
	"fmt"
	"sync"
	"time"
)

type State int

const (
	Alive State = iota
	Dead
	Pending
	Suspect
)

type Identifier string

type Node struct {
	IpAddr         string
	Port           string
	Status         State
	MaxRetry       int
	LastConnected  time.Time
	InitialConnect time.Time
}

type Pulser interface {
	StartPulse(ipAddr, port string) error
	StopPulse()
	AddPulser(ipAddr, port string, maxRetry int) error
	RemovePulser(ipAddr, port string) error
	StopAllPulser()
	StartAllPulser() error
}

type Pulse struct {
	ID       string
	name     string
	status   State
	mutex    sync.Mutex
	nodeMap  map[Identifier]*Node
	detector chan interface{}
}

func Initialize(Capacity int) (*Pulse, error) {
	//TODO
	p := &Pulse{
		ID:       "id",
		name:     "name",
		status:   Alive,
		nodeMap:  make(map[Identifier]*Node),
		detector: make(chan interface{}),
	}

	return p, nil
}

// API's to send heartbeat signals

func CreateNode(ipAddr, port string, maxRetry int) *Node {

	n := &Node{
		IpAddr:   ipAddr,
		Port:     port,
		Status:   Pending,
		MaxRetry: maxRetry,
	}
	return n
}

func AddrToIdentifier(ipAddr, port string) Identifier {
	iden := Identifier(ipAddr + ":" + port)
	return iden
}

// Start responding to pulse messages
func (p *Pulse) StartPulse(ipAddr, port string) error {
	return nil
}

func (p *Pulse) StopPulse() {

}

// API's to detect failures
func (p *Pulse) AddPulser(ipAddr, port string, maxRetry int) error {
	iden := AddrToIdentifier(ipAddr, port)

	if _, ok := p.nodeMap[iden]; !ok {
		return fmt.Errorf("IP Addr and Port combination already exists")
	}
	newNode := CreateNode(ipAddr, port, maxRetry)
	p.mutex.Lock()
	p.nodeMap[iden] = newNode
	p.mutex.Unlock()
	return nil
}

func (p *Pulse) RemovePulser() {

}

func (p *Pulse) StopAllPulser() {

}
