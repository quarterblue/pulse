package pkg

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

var MAGIC_BYTES = []byte("qbpulse")

const initialRTT = 3

type PulseResponse struct {
	Id       uint8
	Addr     string
	Message  string
	Optional interface{}
}

// Create Pulse server
func PulseServer(ctx context.Context, addr Identifier, wg sync.WaitGroup) (net.Addr, error) {
	s, err := net.ListenPacket("udp", string(addr))
	if err != nil {
		return nil, fmt.Errorf("binding udp error %s, %w", addr, err)
	}

	go func() {
		go func() {
			defer wg.Done()
			<-ctx.Done()
			_ = s.Close()
		}()
		rand.Seed(time.Now().UnixNano())

		var idCount uint8 = 0
		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			var writeBuf bytes.Buffer
			encoder := gob.NewEncoder(&writeBuf)
			idCount++
			resp := &PulseResponse{
				Id:      idCount,
				Addr:    s.LocalAddr().String(),
				Message: "Hey there",
			}

			err = encoder.Encode(resp)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("From Coord: %s\n", buf[:n])
			delay := rand.Intn(1000)
			log.Printf("Simulating delay: %d milliseconds\n", delay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
			_, err = s.WriteTo(writeBuf.Bytes(), clientAddr)
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

func SendPulse(ctx context.Context, node *Node, wg sync.WaitGroup) {
	go func(ipAddr, port string) {
		var failedAttempts uint8 = 0
		var avgRTT time.Duration = initialRTT

		log.Printf("Failed Attempts: %d", failedAttempts)
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
				time.Sleep(time.Duration(node.Delay) * time.Second)
				_, err := client.Write(ping)
				start := time.Now()

				if err != nil {
					log.Println("Write failed")
				}

				buf := make([]byte, 1024)
				n, err := client.Read(buf)

				if err != nil {
					failedAttempts++
					log.Printf("Failed Attempt: %d\n", failedAttempts)
					node.mu.RLock()
					if failedAttempts >= node.MaxRetry {
						log.Printf("Failure detected, removing: %s:%s\n", node.IpAddr, node.Port)
						node.mu.RUnlock()
						return
					}
					node.mu.RUnlock()
					continue
				}
				elapsed := time.Since(start)

				if avgRTT == initialRTT {
					log.Println("Initial RTT")
					avgRTT = elapsed
				} else {
					avgRTT = (avgRTT + elapsed) / 2
				}

				dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
				resp := PulseResponse{}
				dec.Decode(&resp)
				log.Printf("Msg Received: %+v\n", resp)
			}
		}
	}(node.IpAddr, node.Port)
}

func RequestPulse(ctx context.Context) {
	fmt.Println("Requesting")
}

type Server struct {
	Payload []byte
	Retires uint8
}
