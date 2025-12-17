package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/GreenDude5/go-port-scanner/internal/api"
	"github.com/GreenDude5/go-port-scanner/internal/scan"
	"github.com/GreenDude5/go-port-scanner/internal/storage"
)

func main() {
	mode := flag.String("mode", "scan", "Mode to run: scan or server")
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

	switch *mode {
	case "server":
		api.StartServer(db)
	case "scan":
		runScanner(db, *hostname, *startPort, *endPort, *threads)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

//	if err := storage.CreateSchema(db); err != nil {
//		log.Fatalf("Failed to create database schema: %v", err)
//	}
func runScanner(db *sql.DB, hostname string, startPort, endPort, threads int) {
	fmt.Printf("Scanning host: %s from port %d to %d\n", hostname, startPort, endPort)

	portsChan := make(chan int, threads)
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range portsChan {
				isOpen := scan.ScanPort("tcp", hostname, p)

				if isOpen {
					fmt.Printf("Port %d is open\n", p)

					err := storage.SaveResult(db, p, "OPEN")
					if err != nil {
						log.Printf("Saving error: %v on port %d", err, p)
					}
				}
			}
		}()
	}

	go func() {
		for i := startPort; i <= endPort; i++ {
			portsChan <- i
		}
		close(portsChan)
	}()

	wg.Wait()
	fmt.Println("Scanning completed.")
}
