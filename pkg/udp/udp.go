package udp

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
)

var MAGIC_BYTES = []byte("Pulse! (From Pulser!)")

func EchoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding udp error %s, %w", addr, err)
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
			_, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}
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

func SendPulse(conn *net.UDPConn, quit chan struct{}) {
	ping := []byte("ping")
	_, err := conn.Write(ping)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Msg Received: %s", buf[:n])
}

type Server struct {
	Payload []byte
	Retires uint8
}

// func (s Server) ListenAndServe(addr string) error {
// 	conn, err := net.ListenPacket("udp", addr)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	log.Printf("Listening on %s ...\n", conn.LocalAddr())
// 	return s.Serve(conn)
// }

// func (s *Server) Serve(conn net.PacketConn) error {
// 	if conn == nil {
// 		return errors.New("nil connection")
// 	}

// 	for {
// 		buf := make([]byte, 1024)
// 		_, addr, err := conn.ReadFrom(buf)
// 		if err != nil {
// 			return err
// 		}

// 	}
// }
