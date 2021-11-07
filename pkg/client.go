package detector

import "sync"

type Pulser struct {
	ID     string
	name   string
	status string
	mutex  sync.Mutex
}

func Initialize() (*Pulser, chan interface{}, error) {
	//TODO
	ch := make(chan interface{})
	p := &Pulser{
		ID:   "id",
		name: "name",
	}

	return p, ch, nil
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
