//go:build !windows
// +build !windows

package protocol

import (
	"net"
	"time"
)

func dial(sockname string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", sockname, timeout)
}
