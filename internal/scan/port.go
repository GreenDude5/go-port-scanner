package scan

import (
	"fmt"
	"net"
	"time"
)

func ScanPort(protocol, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
