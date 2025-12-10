package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/GreenDude5/go-port-scanner/internal/scan"
)

func main() {
	hostname := flag.String("host", "scanme.nmap.org", "Hostname to scan")
	startPort := flag.Int("start", 1, "Start port number")
	endPort := flag.Int("end", 1024, "End port number")
	threads := flag.Int("threads", 100, "Number of concurrent threads")
	flag.Parse()

	fmt.Printf("Scanning host: %s from port %d to %d\n", *hostname, *startPort, *endPort)

	portsChan := make(chan int, 100)

	var wg sync.WaitGroup
	for portRange := 0; portRange <= *threads; portRange++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portsChan {
				isOpen := scan.ScanPort("tcp", *hostname, port)
				if isOpen {
					fmt.Printf("Port %d is open\n", port)
				}
			}
		}()
	}

	for i := *startPort; i <= *endPort; i++ {
		portsChan <- i
	}
	close(portsChan)

	wg.Wait()
	fmt.Println("Scanning completed.")
}
