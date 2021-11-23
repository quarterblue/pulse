package pulse

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
	// Received a pulse directly or indirectly within a time bound
	Alive State = iota
	// Received no pulse directly within a time bound or received gossip message declaring dead
	Dead
	// State assumed when a node is first discovered and contacted
	Pending
	// Suspect state, awaiting time bound to officially announce as dead
	Suspect
)

type Identifier string

// Pulser is a node that responds to pulse requests
type Pulser interface {
	// Starts responding to the pulse requests on IP Address: ipAddr and Port: port
	StartPulseRes(ipAddr, port string) error

	// Stops responding to pulse requests
	StopPulseRes()
}

// Coordinator is a node that requests pulse response from a map of nodes
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

// Pulser and Coordinator functions are indepedent, but can also be used together as a Full node
type FullNode interface {
	Pulser
	Coordinator
	Gossip()
}

type Node struct {
	IpAddr         string
	Port           string
	Status         State
	Tracking       bool
	MaxRetry       uint8
	Delay          uint8
	RTT            float32
	InitialConnect time.Time
	LastConnected  time.Time
	mu             sync.RWMutex
}

// Pulse implements the FullNode interface,
// It can act as a Coordinator or Pulser indenpendently,
// Or it can act as a Full Node with gossip protocol
type Pulse struct {
	Id           string
	name         string
	status       State
	stop         chan interface{}
	mutex        sync.RWMutex
	nodeMap      map[Identifier]*Node
	deadNode     map[Identifier]*Node
	notifyStream chan FailureMessage
	cancelRes    context.CancelFunc
}

type Status struct {
	Id            Identifier
	State         State
	RTT           float32
	LastConnected time.Time
}

func AddrToIdentifier(ipAddr, port string) Identifier {
	return Identifier(ipAddr + ":" + port)
}

func IdentifierToAddr(iden Identifier) (ipAddr, port string) {
	s := strings.SplitN(string(iden), ":", 2)
	return s[0], s[1]
}

func Initialize(capacity int) (*Pulse, chan FailureMessage, error) {

	p := &Pulse{
		Id:           "id",
		name:         "name",
		status:       Alive,
		stop:         make(chan interface{}),
		mutex:        sync.RWMutex{},
		nodeMap:      make(map[Identifier]*Node),
		deadNode:     make(map[Identifier]*Node),
		notifyStream: make(chan FailureMessage, capacity),
	}

	return p, p.notifyStream, nil
}

// API's to send heartbeat signals
func CreateNode(ipAddr, port string, maxRetry, delay uint8) *Node {

	n := &Node{
		IpAddr:   ipAddr,
		Port:     port,
		Status:   Pending,
		Tracking: true,
		MaxRetry: maxRetry,
		Delay:    delay,
		RTT:      3,
		mu:       sync.RWMutex{},
	}

	return n
}

// Start responding to pulse messages
func (p *Pulse) StartPulseRes(ctx context.Context, cancel context.CancelFunc, ipAddr, port string) error {

	var wg sync.WaitGroup
	wg.Add(1)
	addr := AddrToIdentifier(ipAddr, port)
	pulserAddr, err := PulseServer(ctx, addr, wg)
	log.Println("Listening on: ", pulserAddr)

	if err != nil {
		return err
	}

	p.mutex.Lock()
	p.cancelRes = cancel
	p.mutex.Unlock()

	return nil
}

func (p *Pulse) StopPulseRes() {
	p.mutex.Lock()
	p.cancelRes()
	p.mutex.Unlock()
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
	SendPulse(ctx, p.nodeMap[iden], p.notifyStream, wg)
	return nil
}

// Remove the pulser to the map of nodes to monitor and immediately stop sending pulses
func (p *Pulse) RemovePulser(ipAddr, port string) error {
	return nil
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
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	allStatus := make([]*Status, len(p.nodeMap))

	var index int = 0

	for k, v := range p.nodeMap {
		allStatus[index] = &Status{
			Id:            k,
			State:         v.Status,
			RTT:           v.RTT,
			LastConnected: v.LastConnected,
		}
		index++
	}
	return allStatus, nil
}

func (p *Pulse) Gossip() {
	fmt.Println("Start Gossip Protocol")
}
