package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, ports <-chan int, target string) {
	defer wg.Done()

	for port := range ports {
		address := fmt.Sprintf("%s:%d", target, port)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			fmt.Printf("Port %d OPEN\n", port)
			conn.Close()
		}
	}
}

func main() {
	target := "127.0.0.1"
	ports := make(chan int, 100)

	var wg sync.WaitGroup

	// 100 goroutine untuk scanning port, tidak terlalu banyak agar tidak membebani sistem
	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, ports, target)
	}

	for port := 1; port <= 9000; port++ {
		ports <- port
	}
	close(ports)

	wg.Wait()
}