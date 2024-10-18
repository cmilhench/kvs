package kvs

import (
	_ "embed" // The "embed" package must be imported when using go:embed
	"fmt"

	"github.com/cmilhench/kvs/internal/ports"
)

//go:generate sh -c "./scripts/version.sh > .version"
//go:embed .version
var Revision string

type RPCServer struct {
	store ports.Store
}

type GetInput struct {
	Key string
}

type PutAppendInput struct {
	Key   string
	Value string
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
	e := s.store.Put(input.Key, input.Value)
	if e != nil {
		return fmt.Errorf("failed to put, %v", e)
	}
	output.Value = ""
	return nil
}

func (s *RPCServer) Append(input PutAppendInput, output *Output) error {
	old, e := s.store.Append(input.Key, input.Value)
	if e != nil {
		return fmt.Errorf("failed to append, %v", e)
	}
	output.Value = old
	return nil
}

func (s *RPCServer) Get(input GetInput, output *Output) error {
	out, e := s.store.Get(input.Key)
	if e != nil {
		return fmt.Errorf("failed to get, %v", e)
	}
	output.Value = out
	return nil
}
