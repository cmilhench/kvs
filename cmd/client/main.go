package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	service := "RPCServer"
	c, e := rpc.DialHTTP("tcp", ":4444")
	if e != nil {
		log.Fatalf("failed to dial, %v\n", e)
	}

	next := SequenceNumber()
	var reply Reply

	// Testing Put
	arg := Args{Key: "Hello", Value: "World!", Client: 1, Sequence: next()}
	reply = Reply{}
	e = c.Call(fmt.Sprintf("%s.Put", service), arg, &reply)
	if e != nil {
		log.Fatalf("failed to put %+v, %s\n", arg, e)
	}
	if reply.Value != "" {
		log.Fatalf("expected empty reply, got %s\n", reply.Value)
	}

	// Testing Get
	key := Arg{Key: "Hello", Client: 1, Sequence: next()}
	reply = Reply{}
	e = c.Call(fmt.Sprintf("%s.Get", service), key, &reply)
	if e != nil {
		log.Fatalf("failed to get %v, %s\n", key, e)
	}
	if reply.Value != "World!" {
		log.Fatalf("expected World!, got %s\n", reply.Value)
	}

	// Testing Append
	reply = Reply{}
	arg = Args{Key: "Hello", Value: "There!", Client: 1, Sequence: next()}
	e = c.Call(fmt.Sprintf("%s.Append", service), arg, &reply)
	if e != nil {
		log.Fatalf("failed to put %+v, %s\n", arg, e)
	}
	if reply.Value != "World!" {
		log.Fatalf("expected World!, got %s\n", reply.Value)
	}

	// Testing Get
	reply = Reply{}
	key = Arg{Key: "Hello", Client: 1, Sequence: next()}
	e = c.Call(fmt.Sprintf("%s.Get", service), key, &reply)
	if e != nil {
		log.Fatalf("failed to get %v, %s\n", key, e)
	}
	if reply.Value != "World!There!" {
		log.Fatalf("expected World!There!, got %s\n", reply.Value)
	}

	// Testing out of sequence
	reply = Reply{}
	key = Arg{Key: "Hello", Client: 1, Sequence: 4}
	e = c.Call(fmt.Sprintf("%s.Get", service), key, &reply)
	if e == nil {
		log.Fatalf("failed to get out of sequence %v, %s\n", key, e)
	}

	log.Println("Done!")
}

type Arg struct {
	Key      string
	Client   int64
	Sequence int64
}

type Args struct {
	Key      string
	Value    string
	Client   int64
	Sequence int64
}

type Reply struct {
	Value string
}

func SequenceNumber() func() int64 {
	var seq int64
	return func() int64 {
		seq = seq + 1
		return seq
	}
}
