package pkg

import (
	"sync"
	"time"
)

type State int

const (
	Alive   State = iota
	Dead    State = iota
	Pending State = iota
)

type Node struct {
	IpAddr         string
	Port           string
	Status         State
	LastConnected  time.Time
	InitialConnect time.Time
}

type Pulser interface {
	StartPulse(port string) error
	StopPulse()
	AddPulser(node Node, maxRetry int) error
	RemovePulser()
	StopAllPulser()
}

type Pulse struct {
	ID       string
	name     string
	status   State
	mutex    sync.Mutex
	nodeList []*Node
	detector chan interface{}
}

func Initialize(Capacity int) (*Pulse, error) {
	//TODO
	ch := make(chan interface{})
	p := &Pulse{
		ID:       "id",
		name:     "name",
		status:   Alive,
		detector: ch,
	}

	return p, nil
}

// API's to send heartbeat signals

// Start responding to pulse messages
func (p *Pulse) StartPulse(port string) error {
	return nil
}

func (p *Pulse) StopPulse() {

}

// API's to detect failures
func (p *Pulse) AddPulser(node Node, maxRetry int) error {

	return nil
}

func (p *Pulse) RemovePulser() {

}

func (p *Pulse) StopAllPulser() {

}
