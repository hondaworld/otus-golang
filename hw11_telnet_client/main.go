package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Сетевой таймаут")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Использование: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Printf("Ошибка подключения: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
			close(done)
		}
	}()

	go func() {
		_ = client.Receive()
		close(done)
	}()

	select {
	case <-done: // completed sending or receiving
		fmt.Println("Exiting due to EOF or send/receive completion.")
	case <-sigCh: // caught signal
	}

	client.Close()
}
