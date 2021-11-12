package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/quarterblue/pulse/pkg/udp"
)

type Cli struct{}

func (c *Cli) Run() {
	coord := flag.Bool("coord", false, "To start a coordinator")
	pulser := flag.Bool("pulser", false, "To start a pulser")
	portListen := flag.String("port", "9001", "Port to listen on")

	flag.Parse()

	if !(*coord || *pulser) {
		fmt.Println("You must select between a coordinator or a pulser.")
		os.Exit(1)
	}
	fmt.Println("Welcome to Pulse!")
	ctx, cancel := context.WithCancel(context.Background())

	if *coord {
		fmt.Println("Coordinator selected")
		var wg sync.WaitGroup
		ipAddr := "127.0.0.1"
		port := "9005"
		port2 := "9006"
		wg.Add(1)
		udp.SendPulse(ctx, ipAddr, port, wg)
		wg.Add(1)
		udp.SendPulse(ctx, ipAddr, port2, wg)
		wg.Wait()
	} else {
		fmt.Println("Pulser selected")
		var wg sync.WaitGroup
		wg.Add(1)
		ipAddr := "127.0.0.1"
		addr := ipAddr + ":" + *portListen
		pulserAddr, err := udp.PulseServer(ctx, addr, wg)
		fmt.Println("Listening on: ", pulserAddr)
		if err != nil {
			log.Fatal(err)
		}
		wg.Wait()
	}
	// dst, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")

	defer cancel()
}
