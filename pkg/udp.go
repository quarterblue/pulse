package pkg

import (
	"context"
	"fmt"
	"net"
)

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
// 			_, err = s.WriteTo()
// 		}
// 	}()
// 	return s.LocalAddr(), nil
// }
