
<p align="left">
        <img width="20%" src="https://raw.githubusercontent.com/quarterblue/pulse/main/static/pulselogo.png?token=ANKI23IIVVHEVMOXEEE4OYTBUP3DY" alt="Parsec logo">
</p>

---
<a href="https://github.com/quarterblue/pulse/actions/workflows/go.yml" target="_blank">
  <img src="https://github.com/quarterblue/pulse/actions/workflows/go.yml/badge.svg" alt="GitHub Passing">
</a>
<a href="https://github.com/quarterblue/pulse/blob/main/LICENSE" target="_blank">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
</a>

## About

Pulse is an easy-to-use hybrid failure detection library based on simple heartbeat message exchanges overlayed on a gossip protocol. Failure detectors were proposed by <a href="https://dl.acm.org/doi/10.1145/226643.226647">Chandra and Toueg</a> used to solve consensus in asynchronous systems with crash failures. In a fully asynchronous system, a failure detector is impossible to operate. But with time bounds (RTT) we can reasonably suspect a crashed node as failed. The simplest way to build a failure detector would be send and receive heartbeat messages among all nodes in the network.

The problem with heartbeat-based FD is that it is not scalable. Every node in the network exchanges heartbeat message with other nodes, causing the network load to reach an order of O(n^2). For small number of nodes, <= 100, this is a perfectly acceptable way of communicating. However, as the numbers begin to escalate, >= 1000 we are exchanging 1,000,000+ messages! This is where gossip protocol helps us reduce the network load to an order of O(n). In a gossip protocol, every node chooses a random node to (gossip) exchange message with and piggybacks status of other nodes it knows about <a href="https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf">(SWIM)</a>. 

Pulse uses a simple heartbeat protocol when the number of nodes involved are small (<= 100). As the number of nodes grows (customizable by the user), the nodes start disseminating gossip style messages to relay their liveliness. An individual node can opt to keep a heartbeat protocol to receive RTT bounded updates for nodes of their choosing, but the rest of the node discovery will be done via gossip message exchange.

## Features

- Easy to use
- Minimalistic & simple architecture
- Heartbeat sensors (Pulses)
- Based on Gossip protocol
- Dynamic RTT calculation
- Eventually perfect weakly consistent FD
- Easily customizable
- REST API for status updates


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

## References

- <a href="https://dl.acm.org/doi/10.1145/226643.226647">Unreliable failure detectors for reliable distributed systems</a>
- <a href="https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf">SWIM: Scalable Weakly-consistent Infection-style Process Group Membership Protocol</a>
- <a href="https://www.cs.yale.edu/homes/aspnes/pinewiki/FailureDetectors.html">Failure Detectors</a>


## License

Licensed under the MIT License.
