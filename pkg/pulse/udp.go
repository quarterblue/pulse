package pulse

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

var (
	MAGIC_BYTES               = []byte("qbpulse") // Magic bytes allow clients to verify they are speaking the same protocol
	initialRTT  time.Duration = 3 * time.Second   // Intial RTT is generous at the start, but will soon be dynamically calculated
)

// Response struct to send upon receiving a pulse request
type PulseResponse struct {
	Id       uint8
	Addr     string
	Message  string
	Optional interface{}
}

// Message struct to send to notify channel when a node is determined to have failed
type FailureMessage struct {
	Id             Identifier
	State          State
	RetryAttempts  uint8
	GossipPulse    []Identifier
	InitialConnect time.Time
	LastConnected  time.Time
}

// Create a Pulse server that listens on addr Identifier and responds to pulse messages
// Spins up two go routines, one for listening for UDP messages, and the other to wait on ctx.Done
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
				Message: "I am alive!",
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

// Request a pulse signal to a node, if failure is suspected, choose 3 random nodes to inquire about liveliness
// Retry maxRetry times while waiting for average RTT between each messages.
// If both these methods fail to produce pulse from the suspect, mark the node as dead and send a Failure message to nStream channel.
func SendPulse(ctx context.Context, node *Node, nStream chan FailureMessage, wg sync.WaitGroup) {
	go func(ipAddr, port string) {
		var failedAttempts uint8 = 0
		var avgRTT time.Duration = initialRTT

		log.Printf("Failed Attempts: %d", failedAttempts)
		addr := ipAddr + ":" + port

		// dst, err := net.ResolveUDPAddr("udp", addr)

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

				// Pulse response failed, retry request for pulse
				// and immediately activate gossip protcol to ask other nodes for pulse status of the node in question
				if err != nil {
					failedAttempts++
					log.Printf("Failed Attempt: %d\n", failedAttempts)
					node.mu.Lock()
					node.Status = Suspect
					if failedAttempts >= node.MaxRetry {
						defer node.mu.Unlock()
						log.Printf("Failure detected, removing: %s:%s\n", node.IpAddr, node.Port)

						fMsg := FailureMessage{
							Id:             AddrToIdentifier(node.IpAddr, node.Port),
							State:          Dead,
							RetryAttempts:  failedAttempts,
							GossipPulse:    []Identifier{},
							InitialConnect: node.InitialConnect,
							LastConnected:  node.LastConnected,
						}

						nStream <- fMsg
						return
					}
					node.mu.Unlock()
					continue
				}
				recvTime := time.Now()
				elapsed := time.Since(start)

				if avgRTT == initialRTT {
					log.Println("Initial RTT")
					node.mu.Lock()
					node.Status = Alive
					node.mu.Unlock()
					avgRTT = elapsed
				} else {
					avgRTT = (avgRTT + elapsed) / 2
					node.mu.Lock()
					node.LastConnected = recvTime
					node.mu.Unlock()
				}

				dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
				resp := PulseResponse{}
				dec.Decode(&resp)
				log.Printf("Msg Received: %+v\n", resp)
			}
		}
	}(node.IpAddr, node.Port)
}

func RequestPulse(ctx context.Context, node *Node, resp chan interface{}) {
	fmt.Println("Requesting")
}

func PickThree() []*Node {
	return nil
}
