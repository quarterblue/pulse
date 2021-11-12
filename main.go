package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/quarterblue/pulse/pkg/udp"
)

func main() {
	coord := flag.Bool("coord", false, "To start a coordinator")
	pulser := flag.Bool("pulser", false, "To start a pulser")

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
		// addr := "127.0.0.1:9002"

		// s, err := net.ListenPacket("udp", addr)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		wg.Add(1)
		go func(ipAddr, port string) {
			addr := ipAddr + ":" + port
			// dst, err := net.ResolveUDPAddr("udp", addr)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			client, err := net.Dial("udp", addr)
			if err != nil {
				log.Fatal(err)
			}

			ping := []byte("ping")
			for {
				select {
				case <-ctx.Done():
					defer wg.Done()
					return
				default:
					time.Sleep(2 * time.Second)
					_, err := client.Write(ping)
					if err != nil {
						log.Fatal(err)
					}

					buf := make([]byte, 1024)
					n, err := client.Read(buf)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("Msg Received: %s\n", buf[:n])
				}
			}
		}("127.0.0.1", "9005")

		wg.Wait()
	} else {
		fmt.Println("Pulser selected")
		var wg sync.WaitGroup
		wg.Add(1)
		pulserAddr, err := udp.PulseServer(ctx, "127.0.0.1:9005", wg)
		fmt.Println("Listening on: ", pulserAddr)
		if err != nil {
			log.Fatal(err)
		}
		wg.Wait()
	}

	// serverAddr, err := udp.EchoServerUDP(ctx, "127.0.0.1:")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(serverAddr)
	// fmt.Println(reflect.TypeOf(serverAddr))

	// dst, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")

	defer cancel()

	// client, err := net.ListenPacket("udp", "127.0.0.1:")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() {
	// 	_ = client.Close()
	// }()

	// msg := []byte("From pulse!")
	// _, err = client.WriteTo(msg, serverAddr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// buf := make([]byte, 1024)
	// n, addr, err := client.ReadFrom(buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("From: ", addr.String())

	// fmt.Println(string(buf[:n]))
}
