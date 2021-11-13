package pkg

import (
	"context"
	"fmt"
	"strings"
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
	MaxRetry       uint8
	Delay          uint8
	LastConnected  time.Time
	InitialConnect time.Time
}

type Pulser interface {
	// Starts responding the pulse message on IP Address: ipAddr and Port: port
	StartPulseRes(ipAddr, port string) error

	// Stops responding to pulse message
	StopPulseRes()

	// Add an IP Address: ipAddr, Port: port to the monitor list and start asking for pulses
	AddPulser(ipAddr, port string, maxRetry int) error

	// Remove IP Address: ipAddr, Port: port from the monitor list and stop asking for pulses
	RemovePulser(ipAddr, port string) error

	// Collectively stop monitoring all pulsers
	StopAllPulser()

	// Collectively start all pulsers added to monitor list
	StartAllPulser() error

	// Get the current status of a specific Node identified by Identifier
	Status(id Identifier)

	// Get the current status of all Nodes
	StatusAll()
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
func CreateNode(ipAddr, port string, maxRetry, delay uint8) *Node {

	n := &Node{
		IpAddr:   ipAddr,
		Port:     port,
		Status:   Pending,
		MaxRetry: maxRetry,
		Delay:    delay,
	}
	return n
}

func AddrToIdentifier(ipAddr, port string) Identifier {
	return Identifier(ipAddr + ":" + port)
}

func IdentifierToAddr(iden Identifier) (ipAddr, port string) {
	s := strings.SplitN(string(iden), ":", 2)
	return s[0], s[1]
}

// Start responding to pulse messages
func (p *Pulse) StartPulseRes(ipAddr, port string) error {
	return nil
}

func (p *Pulse) StopPulseRes() {

}

// API's to detect failures

// Add the pulser to the map of nodes to monitor and immediately start sending pulses
// maxRetry: the number of times to re-send the udp message before declaring the node dead
// delay: the number of seconds to delay between each message
func (p *Pulse) AddPulser(ipAddr, port string, maxRetry, delay uint8, wg sync.WaitGroup) error {
	iden := AddrToIdentifier(ipAddr, port)

	if _, ok := p.nodeMap[iden]; ok {
		return fmt.Errorf("IP Addr and Port combination already exists")
	}
	newNode := CreateNode(ipAddr, port, maxRetry, delay)
	p.mutex.Lock()
	p.nodeMap[iden] = newNode
	p.mutex.Unlock()

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	SendPulse(ctx, p.nodeMap[iden], wg)
	return nil
}

func (p *Pulse) RemovePulser() {

}

func (p *Pulse) StopAllPulser() {

}
