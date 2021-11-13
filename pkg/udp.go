package pkg

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var MAGIC_BYTES = []byte("Pulse! (From Pulser!)")

// Create Pulse server
func PulseServer(ctx context.Context, addr string, wg sync.WaitGroup) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding udp error %s, %w", addr, err)
	}

	go func() {
		go func() {
			defer wg.Done()
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}
			fmt.Printf("From Coord: %s", buf[:n])
			_, err = s.WriteTo(MAGIC_BYTES, clientAddr)
			if err != nil {
				return
			}
		}
	}()
	return s.LocalAddr(), nil
}

// func ListenPulse(conn *net.UDPConn, quit chan struct{}) {
// 	buf := make([]byte, 1024)
// 	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
// 	for err == nil {
// 		n, remoteAddr, err = conn.ReadFromUDP(buf)
// 		fmt.Println("from", remoteAddr)
// 	}
// }

func SendPulse(ctx context.Context, n *Node, wg sync.WaitGroup) {
	go func(ipAddr, port string) {
		var failedAttempts uint8 = 0
		fmt.Println(failedAttempts)
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
				time.Sleep(time.Duration(n.Delay) * time.Second)
				_, err := client.Write(ping)
				if err != nil {
					failedAttempts++
					if failedAttempts >= n.MaxRetry {
						log.Fatal(err)
					}
				}

				buf := make([]byte, 1024)
				n, err := client.Read(buf)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Msg Received: %s\n", buf[:n])
			}
		}
	}(n.IpAddr, n.Port)
}

type Server struct {
	Payload []byte
	Retires uint8
}
