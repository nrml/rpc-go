package rpc

import (
	"fmt"
	"log"
	//"net/rpc"
	"testing"
)

type Registration struct {
	ID       int64
	Email    string
	Password string
}

type Status struct {
	Namespace string
	Status    int8
}

var (
	address = "127.0.0.1"
	port    = 9001
	remote  = fmt.Sprintf("%s:%d", address, port)
)

func Test_Server(t *testing.T) {

	client, err := NewClient("testsecret", "testnamespace", "Membership")
	err = client.Connect(address, int64(port))

	if err != nil {
		log.Fatal("dialing:", err)
		return
	}

	reg := new(Registration)
	//synchronous call
	err = client.Call("BadFunc", reg)

	if err != nil {
		log.Println("good that service return error with bad func", err.Error())
	}

	//client conn shuts down, reconnect

	err = client.Connect(address, int64(port))

	send := Registration{0, "dummy9001", "dummyPass"}

	reg = new(Registration)

	async := client.Async("Create", send, reg)

	cb := <-async.Done

	if cb.Error != nil {
		log.Println(cb.Error.Error())
		return
	}

	log.Printf("created reg: %d %s %s\n", reg.ID, reg.Email, reg.Password)

	id := reg.ID
	reg = new(Registration)
	async = client.Async("Get", id, reg)

	cb = <-async.Done

	if cb.Error != nil {
		log.Println(cb.Error.Error())
		return
	}

	log.Printf("got reg: %d %s %s\n", reg.ID, reg.Email, reg.Password)

	reg = new(Registration)

	async = client.Async("Login", send, reg)

	cb = <-async.Done

	if cb.Error != nil {
		log.Println("ERROR: " + cb.Error.Error())
	}

	log.Printf("logged in reg: %d %s %s\n", reg.ID, reg.Email, reg.Password)

	list := new([]Registration)
	async = client.Async("List", list)

	cb = <-async.Done

	if cb.Error != nil {
		log.Println("ERROR: " + cb.Error.Error())
	}

	log.Printf("list of reges: %v\n", list)

	list = new([]Registration)
	async = client.Async("Search", "email='dummy9001'", list)

	cb = <-async.Done

	if cb.Error != nil {
		log.Println("ERROR: " + cb.Error.Error())
	}

	log.Printf("results of search: %v\n", list)

}
