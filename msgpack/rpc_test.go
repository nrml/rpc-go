package rpc

import (
	"fmt"
	"log"
	//"net/rpc"
	"testing"
	"time"
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
	//port    = 20000
	port   = 9002
	remote = fmt.Sprintf("%s:%d", address, port)
)

func NoTest_MulitpleConns(t *testing.T) {
	log.Println("testing multiple connections")
	clt, _ := NewClient(fmt.Sprintf("testkey%d", 1), fmt.Sprintf("namespace%d", 1), "Membership")
	clt.Connect(address, int64(port))

	clt2, _ := NewClient(fmt.Sprintf("testkey%d", 2), fmt.Sprintf("namespace%d", 1), "Membership")
	clt2.Connect(address, int64(port))

	clt3, _ := NewClient(fmt.Sprintf("testkey%d", 3), fmt.Sprintf("namespace%d", 1), "Membership")
	clt3.Connect(address, int64(port))

	for i := 1; i < 200; i++ {
		go func(idx int) {
			st := new(status)
			clt.Call("Init", clt.key, clt.namespace, st)
		}(i)
		go func(idx int) {
			st := new(status)
			clt2.Call("Init", clt2.key, clt.namespace, st)
		}(i)
		go func(idx int) {
			st := new(status)
			clt3.Call("Init", clt3.key, clt.namespace, st)
		}(i)
	}
	log.Println("sleeping")
	time.Sleep(10 * time.Second)
	log.Println("done")
}

func Test_Server(t *testing.T) {
	log.Println("begin server test")
	client, err := NewClient("testsecret", "testnamespace", "Membership")

	err = client.Connect(address, int64(port))
	log.Println("checking for dial error")
	if err != nil {
		log.Println("killing server test after dial failure")
		log.Fatal("dialing:", err)
		return
	} else {
		log.Println("no connect error found")
	}
	//return
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
