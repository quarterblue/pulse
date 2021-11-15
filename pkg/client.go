package pkg

import (
	"context"
	"fmt"
	"log"
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
	RTT            float32
	InitialConnect time.Time
	LastConnected  time.Time
	mu             sync.RWMutex
}

type Pulser interface {
	// Starts responding to the pulse requests on IP Address: ipAddr and Port: port
	StartPulseRes(ipAddr, port string) error

	// Stops responding to pulse requests
	StopPulseRes()
}

type Coordinator interface {
	// Add an IP Address: ipAddr, Port: port to the monitor list and start asking for pulses
	AddPulser(ipAddr, port string, maxRetry int) error

	// Remove IP Address: ipAddr, Port: port from the monitor list and stop asking for pulses
	RemovePulser(ipAddr, port string) error

	// Collectively stop monitoring all pulsers
	StopAllPulser()

	// Collectively start all pulsers added to monitor list
	StartAllPulser() error

	// Get the current status of a specific Node identified by Identifier
	Status(id Identifier) (Status, error)

	// Get the current status of all Nodes
	StatusAll() ([]*Status, error)
}

type FullNode interface {
	Pulser
	Coordinator
}

type Pulse struct {
	ID           string
	name         string
	status       State
	stop         chan interface{}
	mutex        sync.RWMutex
	nodeMap      map[Identifier]*Node
	deadNode     map[Identifier]*Node
	notifyStream chan interface{}
}

type Status struct {
	Id            Identifier
	State         State
	RTT           float32
	LastConnected time.Time
}

func Initialize(Capacity int) (*Pulse, error) {
	//TODO
	p := &Pulse{
		ID:           "id",
		name:         "name",
		status:       Alive,
		stop:         make(chan interface{}),
		mutex:        sync.RWMutex{},
		nodeMap:      make(map[Identifier]*Node),
		deadNode:     make(map[Identifier]*Node),
		notifyStream: make(chan interface{}),
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
		RTT:      3,
		mu:       sync.RWMutex{},
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
func (p *Pulse) StartPulseRes(ctx context.Context, ipAddr, port string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	addr := AddrToIdentifier(ipAddr, port)
	pulserAddr, err := PulseServer(ctx, addr, wg)
	log.Println("Listening on: ", pulserAddr)

	if err != nil {
		return err
	}

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
	p.mutex.RLock()
	if _, ok := p.nodeMap[iden]; ok {
		return fmt.Errorf("IP Addr and Port combination already exists")
	}
	p.mutex.RUnlock()
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

// Collectively start all pulsers added to monitor list
func (p *Pulse) StartAllPulser() error {
	return nil
}

// Get the current status of a specific Node identified by Identifier
func (p *Pulse) Status(id Identifier) (*Status, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if node, ok := p.nodeMap[id]; !ok {
		return nil, fmt.Errorf("IP Addr and Port combination does not exists")
	} else {
		s := &Status{
			Id:            id,
			State:         node.Status,
			RTT:           node.RTT,
			LastConnected: node.LastConnected,
		}
		return s, nil
	}
}

// Get the current status of all Nodes
func (p *Pulse) StatusAll() ([]*Status, error) {
	return nil, nil
}
