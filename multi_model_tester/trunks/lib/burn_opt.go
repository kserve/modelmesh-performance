package trunks

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// BurnOpt ...
type BurnOpt func(*Burner)

// NumWorker : Deprecating; use WithNumWorker(uint64)
func NumWorker(num uint64) BurnOpt {
	return func(b *Burner) { b.numWorker = num }
}

// WithLoop : Deprecating; use WithLooping(bool)
func WithLoop() BurnOpt {
	return func(b *Burner) { b.loop = true }
}

// WithNumWorker ...
func WithNumWorker(num uint64) BurnOpt {
	return func(b *Burner) { b.numWorker = num }
}

// WithNumConnPerHost ...
func WithNumConnPerHost(num uint64) BurnOpt {
	return func(b *Burner) { b.numConnPerHost = num }
}

// WithLooping ...
func WithLooping(yesno bool) BurnOpt {
	return func(b *Burner) { b.loop = yesno }
}

// WithDumpFile ...
func WithDumpFile(fileName string) BurnOpt {
	return func(b *Burner) {
		b.dumpFile = fileName
		if fileName != "" {
			b.dump = true
		}
	}
}

// WithMetadata ...
func WithMetadata(md metadata.MD) BurnOpt {
	return func(b *Burner) {
		b.ctx = metadata.NewOutgoingContext(context.Background(), md)
	}
}

// WithMaxRecvSize ...
func WithMaxRecvSize(s int) BurnOpt {
	return func(b *Burner) {
		b.maxRecv = s
	}
}

// WithMaxSendSize ...
func WithMaxSendSize(s int) BurnOpt {
	return func(b *Burner) {
		b.maxSend = s
	}
}
