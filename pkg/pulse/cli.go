package pulse

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
	restApi := flag.Bool("api", false, "Start REST Api")

	flag.Parse()

	if !(*coord || *pulser) {
		log.Println("You must select between a coordinator or a pulser.")
		os.Exit(1)
	}
	log.Println("Welcome to Pulse!")
	ctx, cancel := context.WithCancel(context.Background())

	if *coord {
		log.Println("Coordinator selected")
		p, nStream, err := Initialize(10)
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		nodeList := []string{"9005", "9006", "9007"}

		for _, node := range nodeList {
			wg.Add(1)
			err = p.AddPulser("127.0.0.1", node, 3, 2, wg)
			if err != nil {
				log.Println(err)
			}
		}

		go func(stream chan FailureMessage) {
			log.Println(<-stream)
		}(nStream)

		if *restApi {
			go HttpAPI(p, 7001, "Development")
		}

		wg.Wait()
	} else {
		log.Println("Pulser selected")
		var wg sync.WaitGroup
		wg.Add(1)
		p, _, err := Initialize(10)
		if err != nil {
			log.Fatal(err)
		}
		err = p.StartPulseRes(ctx, cancel, "127.0.0.1", *portListen)
		if err != nil {
			log.Fatal(err)
		}

		wg.Wait()
	}
	// dst, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")

	defer cancel()
}
