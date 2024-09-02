package tools

import (
	"fmt"
	"net"
	"time"
)

func TestTCPConnection(address string) error {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to TCP service: %w", err)
	}
	conn.Close()
	return nil
}
