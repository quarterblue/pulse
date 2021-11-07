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

type Pulser struct {
	ID       string
	name     string
	status   State
	mutex    sync.Mutex
	nodeList []*Node
	detector chan interface{}
}

func Initialize(Capacity int) (*Pulser, error) {
	//TODO
	ch := make(chan interface{})
	p := &Pulser{
		ID:       "id",
		name:     "name",
		detector: ch,
	}

	return p, nil
}

// API's to send heartbeat signals

// Start responding to pulse messages
func (p *Pulser) StartPulse(port string) error {
	return nil
}

func (p *Pulser) StopPulse() {

}

// API's to detect failures
func (p *Pulser) AddMonitor() {

}

func (p *Pulser) RemoveMonitor() {

}

func (p *Pulser) StopMonitoring() {

}
