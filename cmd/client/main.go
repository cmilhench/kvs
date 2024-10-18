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

	var reply Reply

	// Testing Put
	arg := Args{Key: "Hello", Value: "World!"}
	reply = Reply{}
	e = c.Call(fmt.Sprintf("%s.Put", service), arg, &reply)
	if e != nil {
		log.Fatalf("failed to put %+v, %s\n", arg, e)
	}
	if reply.Value != "" {
		log.Fatalf("expected empty reply, got %s\n", reply.Value)
	}

	// Testing Get
	key := Arg{Key: "Hello"}
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
	arg = Args{Key: "Hello", Value: "There!"}
	e = c.Call(fmt.Sprintf("%s.Append", service), arg, &reply)
	if e != nil {
		log.Fatalf("failed to put %+v, %s\n", arg, e)
	}
	if reply.Value != "World!" {
		log.Fatalf("expected World!, got %s\n", reply.Value)
	}

	// Testing Get
	reply = Reply{}
	key = Arg{Key: "Hello"}
	e = c.Call(fmt.Sprintf("%s.Get", service), key, &reply)
	if e != nil {
		log.Fatalf("failed to get %v, %s\n", key, e)
	}
	if reply.Value != "World!There!" {
		log.Fatalf("expected World!There!, got %s\n", reply.Value)
	}

	log.Println("Done!")
}

type Arg struct {
	Key string
}

type Args struct {
	Key   string
	Value string
}

type Reply struct {
	Value string
}
