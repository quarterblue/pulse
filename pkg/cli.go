package pkg

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
)

type Cli struct {
	options string
}

func (c *Cli) Run() {
	coord := flag.Bool("coord", false, "To start a coordinator")
	pulser := flag.Bool("pulser", false, "To start a pulser")
	portListen := flag.String("port", "9001", "Port to listen on")

	flag.Parse()

	if !(*coord || *pulser) {
		log.Println("You must select between a coordinator or a pulser.")
		os.Exit(1)
	}
	log.Println("Welcome to Pulse!")
	ctx, cancel := context.WithCancel(context.Background())

	if *coord {
		log.Println("Coordinator selected")
		p, err := Initialize(10)
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		err = p.AddPulser("127.0.0.1", "9005", 3, 5, wg)
		if err != nil {
			log.Println(err)
		}
		wg.Add(1)
		err = p.AddPulser("127.0.0.1", "9006", 3, 2, wg)
		if err != nil {
			log.Println(err)
		}
		// ipAddr := "127.0.0.1"
		// port := "9005"
		// port2 := "9006"
		// wg.Add(1)
		// SendPulse(ctx, ipAddr, port, wg)
		// wg.Add(1)
		// SendPulse(ctx, ipAddr, port2, wg)
		wg.Wait()
	} else {
		log.Println("Pulser selected")
		var wg sync.WaitGroup
		wg.Add(1)
		p, err := Initialize(10)
		if err != nil {
			log.Fatal(err)
		}
		err = p.StartPulseRes(ctx, "127.0.0.1", *portListen)
		if err != nil {
			log.Fatal(err)
		}

		wg.Wait()
	}
	// dst, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")

	defer cancel()
}
