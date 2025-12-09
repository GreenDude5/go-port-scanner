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
	flag.Parse()

	fmt.Printf("Scanning host: %s from port %d to %d\n", *hostname, *startPort, *endPort)

	var wg sync.WaitGroup
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			isOpen := scan.ScanPort("tcp", *hostname, p)
			if isOpen {
				fmt.Printf("Port %d is open\n", p)
			}
		}(port)
	}

	wg.Wait()
	fmt.Println("Scanning completed.")
}
