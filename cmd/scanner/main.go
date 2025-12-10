package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/GreenDude5/go-port-scanner/internal/scan"
	"github.com/GreenDude5/go-port-scanner/internal/storage"
)

func main() {
	hostname := flag.String("host", "scanme.nmap.org", "Hostname to scan")
	startPort := flag.Int("start", 1, "Start port number")
	endPort := flag.Int("end", 1024, "End port number")
	threads := flag.Int("threads", 100, "Number of concurrent threads")
	flag.Parse()

	db, err := storage.NewConnection("scanner_user", "scanner_pass", "localhost:5432", "scanner_db")
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	if err := storage.CreateSchema(db); err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}

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
					err := storage.SaveResult(db, port, "open")
					if err != nil {
						log.Printf("Failed to save result for port %d: %v", port, err)
					}
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
