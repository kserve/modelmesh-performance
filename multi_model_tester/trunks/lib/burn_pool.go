package trunks

import (
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

var (
	errPoolNoConn = fmt.Errorf("no connection in pool xxx")
)

// simple client-side round-robin pool
type pool struct {
	conns []*grpc.ClientConn
	mu    sync.Mutex
	next  int
}

func (p *pool) pick() (*grpc.ClientConn, error) {
	size := len(p.conns)
	if size == 0 {
		return nil, errPoolNoConn
	}

	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
		p.next = (p.next + 1) % size
	}()
	return p.conns[p.next], nil
}

func (p *pool) close() error {
	var errs []string
	for _, c := range p.conns {
		if err := c.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}
