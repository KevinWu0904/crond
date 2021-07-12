package raft

import (
	"net"
	"time"

	"github.com/hashicorp/raft"
)

// StreamLayer implements raft low-level network transport.
type StreamLayer struct {
	net.Listener
}

// NewStreamLayer creates StreamLayer.
func NewStreamLayer(listener net.Listener) *StreamLayer {
	return &StreamLayer{
		listener,
	}
}

// Dial connects to the address on the named network.
func (t *StreamLayer) Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	return dialer.Dial("tcp", string(address))
}
