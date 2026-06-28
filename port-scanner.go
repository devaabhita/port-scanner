package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Result struct {
	Port    int
	Service string
}

func detectService(conn net.Conn, port int) string {
	// coba ambil banner
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	banner := string(buf[:n])

	// deteksi sederhana service yang open
	switch {
	case port == 80 || port == 8080:
		return "HTTP"
	case port == 443:
		return "HTTPS"
	case port == 22:
		return "SSH"
	case port == 3306:
		return "MySQL"
	case port == 5432:
		return "PostgreSQL"
	case len(banner) > 0:
		return "Unknown (" + banner[:min(20, len(banner))] + ")"
	default:
		return "Unknown"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func worker(wg *sync.WaitGroup, ports <-chan int, results chan<- Result, target string) {
	defer wg.Done()

	for port := range ports {
		address := net.JoinHostPort(target, fmt.Sprintf("%d", port))

		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			continue
		}

		service := detectService(conn, port)

		results <- Result{
			Port:    port,
			Service: service,
		}

		conn.Close()
	}
}

func main() {
	target := "localhost"

	startPort := 1
	endPort := 9000
	// numworkers atau goroutine yang akan dijalankan tidak banyak agar tidak membebani sistem
	numWorkers := 100

	ports := make(chan int, 100)
	results := make(chan Result)

	var wg sync.WaitGroup

	// start workers/goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, ports, results, target)
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			ports <- port
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("PORT\tSTATUS\tSERVICE")
	for res := range results {
		fmt.Printf("%d\tOPEN\t%s\n", res.Port, res.Service)
	}
}