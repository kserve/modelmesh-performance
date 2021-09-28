package trunks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(codecIgnoreResp{})
}

// Burner is who runs the stress test against Gtargets.
type Burner struct {
	numWorker      uint64
	numConnPerHost uint64
	loop           bool
	dump           bool
	dumpFile       string
	maxRecv        int
	maxSend        int

	pool   *pool
	ctx    context.Context
	stopch chan struct{}
}

// NewBurner creates a new Burner with burner options.
func NewBurner(hosts []string, opts ...BurnOpt) (*Burner, error) {
	if len(hosts) < 1 {
		return nil, fmt.Errorf("no host")
	}

	b := &Burner{
		stopch:    make(chan struct{}),
		numWorker: DefaultWorkers,
		ctx:       context.Background(),
	}
	fmt.Println(&opts)
	for _, opt := range opts {
		opt(b)
	}

	var cos []grpc.CallOption
	if b.maxRecv > 0 {
		cos = append(cos, grpc.MaxCallRecvMsgSize(b.maxRecv))
	}
	if b.maxSend > 0 {
		cos = append(cos, grpc.MaxCallSendMsgSize(b.maxSend))
	}

	p := &pool{}

	for _, h := range hosts {
		var x uint64
		for ; x < b.numConnPerHost; x++ {
			var c *grpc.ClientConn
			var err error

			if len(cos) > 0 {
				c, err = grpc.Dial(h, grpc.WithDefaultCallOptions(cos...), grpc.WithInsecure())
			} else {
				c, err = grpc.Dial(h, grpc.WithInsecure())
			}

			if err != nil {
				return nil, fmt.Errorf("failed to dial %s: %v", h, err)
			}
			p.conns = append(p.conns, c)
		}
	}

	if len(p.conns) < 1 {
		return nil, fmt.Errorf("no connection in pool")
	}

	b.pool = p
	return b, nil
}

// WaitDumpDone waits until all response dumpings are done
// and write everything into the dump file
func (b *Burner) WaitDumpDone() error {
	// dummy-proof
	if !b.dump {
		return nil
	}

	dumpers.Wait()

	f, err := os.OpenFile(b.dumpFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(memBuf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (b *Burner) Close() error {
	return b.pool.close()
}

func (b *Burner) Burn(tgt *Gtarget, rate uint64, du time.Duration) <-chan *Result {
	var workers sync.WaitGroup
	results := make(chan *Result)
	ticks := make(chan time.Time)
	for i := uint64(0); i < b.numWorker; i++ {
		workers.Add(1)
		go b.burn(tgt, &workers, ticks, results)
	}

	go func() {
		defer close(results)
		defer workers.Wait()
		defer close(ticks)
		interval := 1e9 / rate
		hits := rate * uint64(du.Seconds())
		began, done := time.Now(), uint64(0)
		for {
			now, next := time.Now(), began.Add(time.Duration(done*interval))
			time.Sleep(next.Sub(now))
			select {
			case ticks <- max(next, now):
				if done++; done == hits {
					return
				}
			case <-b.stopch:
				return
			default:
				workers.Add(1)
				go b.burn(tgt, &workers, ticks, results)
			}
		}
	}()

	return results
}

func (b *Burner) Stop() {
	select {
	case <-b.stopch:
		return
	default:
		close(b.stopch)
	}
}

func (b *Burner) burn(tgt *Gtarget, workers *sync.WaitGroup, ticks <-chan time.Time, results chan<- *Result) {
	defer workers.Done()
	for tk := range ticks {
		results <- b.hit(tgt, tk)
	}
}

func (b *Burner) hit(tgt *Gtarget, tm time.Time) (res *Result) {
	res = &Result{Timestamp: tm}
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}

		res.Latency = time.Since(tm)

		// count failure by cheating metrics
		// as if this is an HTTP call
		if err != nil {
			if res.Code < 100 {
				// if res.code is not set
				// consider it the fault from caller side
				res.Code = 400
			}
			res.Error = err.Error()
		} else {
			res.Code = 200
		}
	}()

	c, err := b.pool.pick()
	if err != nil {
		b.Stop()
		return
	}

	req, err := tgt.pick(b.loop)
	if err != nil {
		b.Stop()
		return
	}

	if !b.dump {
		// simply discard the response, no operation on that object
		// so it's able to share one single response object
		err = c.Invoke(b.ctx, tgt.MethodName, req, tgt.Response, grpc.CallContentSubtype("proto-ignore-resp"))
	} else {
		// clone a response object to avoid race-condition
		bb := reflect.New(reflect.TypeOf(tgt.Response).Elem()).Interface().(proto.Message)
		err = c.Invoke(b.ctx, tgt.MethodName, req, bb)
		if err == nil {
			// start a new routine to serialize and dump
			// to avoid blocking processing
			dumpers.Add(1)
			go func(resp proto.Message) {
				defer dumpers.Done()
				bytes, err := json.Marshal(resp)
				if err == nil {
					memBufMutex.Lock()
					defer memBufMutex.Unlock()
					bytes = append(bytes, []byte{'\r', '\n'}...)
					memBuf.Write(bytes)
				}
			}(bb)
		}
	}
	return
}
