package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "-timeout 10s")
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("Not enough arguments")
	}
	addr := net.JoinHostPort(flag.Args()[0], flag.Args()[1])
	tc := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)

	if err := tc.Connect(); err != nil {
		log.Fatal(err)
	}
	log.SetFlags(0)
	log.Println("...Connected to", addr)
	defer tc.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		err := tc.Receive()
		if err != nil {
			log.Println("receive:", err)
		} else {
			log.Println("...Connection was closed by peer")
		}
		cancel()
	}()

	go func() {
		err := tc.Send()
		if err != nil {
			log.Println("send:", err)
		} else {
			log.Println("...EOF")
		}
		cancel()
	}()

	<-ctx.Done()
}
