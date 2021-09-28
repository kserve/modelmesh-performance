package trunks

import (
	"errors"
	"sync"

	"github.com/golang/protobuf/proto"
)

var errNoRequest = errors.New("no available request")

// Gtarget represents an attack target i.e one gRPC endpoint
type Gtarget struct {
	// gRPC method name;
	// Setter & Getter
	MethodName string

	// the gRPC request message(s); Could be multiple.
	// This enables `trunks` to read and store multiple
	// requests to burn one enpoint.
	// Setter & Getter
	Requests []proto.Message

	// the gRPC response message
	// Getter
	Response proto.Message

	ipick   uint64 // index of the request to pick next
	ipickMx sync.Mutex
}

func (t *Gtarget) pick(loop bool) (proto.Message, error) {
	nReq := uint64(len(t.Requests))

	if nReq == 0 {
		return nil, errNoRequest
	}

	if nReq == 1 {
		return t.Requests[0], nil
	}

	t.ipickMx.Lock()
	defer t.ipickMx.Unlock()

	if !loop && t.ipick >= nReq {
		return nil, errNoRequest
	}

	defer func() {
		t.ipick++
	}()

	return t.Requests[t.ipick%nReq], nil
}
