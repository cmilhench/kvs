package kvs

import (
	_ "embed" // The "embed" package must be imported when using go:embed
	"fmt"
	"sync"

	"github.com/cmilhench/kvs/internal/ports"
)

//go:generate sh -c "./scripts/version.sh > .version"
//go:embed .version
var Revision string

type RPCServer struct {
	store ports.Store
	mu    sync.Mutex
	table map[int64]int64
}

type GetInput struct {
	Key      string
	Client   int64
	Sequence int64
}

type PutAppendInput struct {
	Key      string
	Value    string
	Client   int64
	Sequence int64
}

type Output struct {
	Value string
}

func NewRPCServer(store ports.Store) *RPCServer {
	return &RPCServer{
		store: store,
		table: make(map[int64]int64),
	}
}

func (s *RPCServer) Put(input PutAppendInput, output *Output) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	if last, ok := s.table[input.Client]; ok {
		if input.Sequence <= last {
			return fmt.Errorf("out of sequence")
		}
	}
	s.table[input.Client] = input.Sequence
	e := s.store.Put(input.Key, input.Value)
	if e != nil {
		return fmt.Errorf("failed to put, %v", e)
	}
	output.Value = ""
	return nil
}

func (s *RPCServer) Append(input PutAppendInput, output *Output) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	if last, ok := s.table[input.Client]; ok {
		if input.Sequence <= last {
			return fmt.Errorf("out of sequence")
		}
	}
	s.table[input.Client] = input.Sequence
	old, e := s.store.Append(input.Key, input.Value)
	if e != nil {
		return fmt.Errorf("failed to append, %v", e)
	}
	output.Value = old
	return nil
}

func (s *RPCServer) Get(input GetInput, output *Output) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()
	if last, ok := s.table[input.Client]; ok {
		if input.Sequence <= last {
			return fmt.Errorf("out of sequence")
		}
	}
	s.table[input.Client] = input.Sequence
	out, e := s.store.Get(input.Key)
	if e != nil {
		return fmt.Errorf("failed to get, %v", e)
	}
	output.Value = out
	return nil
}
