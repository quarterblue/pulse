
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
	StartPulse(port string) error
	StopPulse()
	AddPulser(node Node, maxRetry int) error
	RemovePulser()
	StopAllPulser()
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
