package udp

import (
	"context"
	"fmt"
	"log"
	"net"
)

var MAGIC_BYTES = []byte("pulse")

func EchoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("Binding to UDP %s: %w", addr, err)
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)
		ano := []byte("Greetings from Pulser Go.")
		for {
			_, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			_, err = s.WriteTo(ano[:], clientAddr)
			if err != nil {
				return
			}
		}
	}()
	return s.LocalAddr(), nil
}

// Create Pulse server
// func PulseServer(ctx context.Context, addr string) (net.Addr, error) {
// 	s, err := net.ListenPacket("udp", addr)
// 	if err != nil {
// 		return nil, fmt.Errorf("Binding UDP Error %s, %w", addr, err)
// 	}

// 	go func() {
// 		go func() {
// 			<-ctx.Done()
// 			_ = s.Close()
// 		}()

// 		buf := make([]byte, 1024)
// 		for {
// 			_, clientAddr, err := s.ReadFrom(buf)
// 			if err != nil {
// 				return
// 			}
// 			_, err = s.WriteTo(clientAddr)
// 		}
// 	}()
// 	return s.LocalAddr(), nil
// }

// func ListenPulse(conn *net.UDPConn, quit chan struct{}) {
// 	buf := make([]byte, 1024)
// 	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
// 	for err == nil {
// 		n, remoteAddr, err = conn.ReadFromUDP(buf)
// 		fmt.Println("from", remoteAddr)
// 	}
// }

func SendPulse(conn *net.UDPConn, ipAddr string, port string, quit chan struct{}) {
	addr := ipAddr + ":" + port
	dst, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.WriteTo(MAGIC_BYTES, dst)
	if err != nil {
		log.Fatal(err)
	}
}
