package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/cmilhench/kvs"
	"github.com/cmilhench/kvs/internal/adapters"
)

func main() {
	server := kvs.NewRPCServer(adapters.NewMemoryStore())
	e := rpc.Register(server)
	if e != nil {
		log.Fatalf("failed to register, %v\n", e)
	}
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":4444")
	if e != nil {
		log.Fatalf("failed to listen, %v\n", e)
	}
	log.Printf("Serving RPC server on port %d", 4444)
	e = http.Serve(l, nil)
	if e != nil {
		log.Fatalf("shutdown %v\n", e)
	}

	// TODO: take a variable that clusters this and implement replication
}
