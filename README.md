
<p align="left">
        <img width="20%" src="https://raw.githubusercontent.com/quarterblue/pulse/main/static/pulselogo.png?token=ANKI23LN4KTYDEHVJKQIFVDBSF7L4" alt="Parsec logo">
</p>

---

## About

Easy to use failure detection library based on gossip protocol

## Features

- Easy to use
- Minimalistic, simple architecture
- Heartbeat sensors (Pulses)
- Based on Gossip protocol
- Dynamic RTT calculation
- Best effort FD
- Easily customizable


## Installation

To install in Unix:

```shell
$ cd projectdir/
$ go get github.com/quarterblue/pulse
```

Import into your Go project:

```go
import (
  	"github.com/quarterblue/pulse"
)
```


## Usage

Pulser implementions the following interface:

```go
type Pulser interface {
	// Starts responding to the pulse requests on IP Address: ipAddr and Port: port
	StartPulseRes(ipAddr, port string) error

	// Stops responding to pulse requests
	StopPulseRes()
}
```

Coordinator implementions the following interface:

```go
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
	Status(id Identifier)

	// Get the current status of all Nodes
	StatusAll()
}
```

To initialize and add a Pulser

```go
package main

import (
        "log"
        "github.com/quarterblue/pulse"
)

func main() {
        capacitySize := 10
        p, err := pulse.Initialize(capacitySize)
        if err != nil {
                log.Fatal(err)
        }
        
        err = p.AddPulser(node, 3)
        if err != nil {
                log.Fatal(err)
        }
        
}
```

## License

Licensed under the MIT License.
