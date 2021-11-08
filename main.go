package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"

	"github.com/quarterblue/pulse/pkg/udp"
)

func main() {
	cl := flag.Bool("client", false, "To start a client")
	ser := flag.Bool("server", false, "To start a server")

	flag.Parse()

	if !(*cl || *ser) {
		fmt.Println("You must select between a client or a server.")
		os.Exit(1)
	}
	if *cl {
		fmt.Println("Client selected")
	} else {
		fmt.Println("Server selected")
	}

	fmt.Println("Welcome to Pulse!")
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := udp.EchoServerUDP(ctx, "127.0.0.1:")
	fmt.Println(serverAddr)
	fmt.Println(reflect.TypeOf(serverAddr))

	dst, err := net.ResolveUDPAddr("udp", "127.0.0.1:9001")

	fmt.Println("h")
	fmt.Println(dst)

	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	msg := []byte("From pulse!")
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("From: ", addr.String())

	fmt.Println(string(buf[:n]))
}
