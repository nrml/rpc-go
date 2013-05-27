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
	}

	create := Registration{0, "dummy9001", "dummyPass"}

	reg := new(Registration)

	// Synchronous call
	err = client.Call("Create", create, reg)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("created reg: %d %s %s\n", reg.ID, reg.Email, reg.Password)

	get := new(Registration)
	err = client.Call("Get", reg.ID, get)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("got reg: %d %s %s\n", get.ID, get.Email, get.Password)

	login := new(Registration)
	logreq := Registration{}
	logreq.Email = reg.Email
	logreq.Password = "dummyPass"
	err = client.Call("Login", logreq, login)

	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	log.Printf("logged in reg: %d %s %s\n", login.ID, login.Email, login.Password)

	regs := new([]Registration)
	err = client.Call("List", regs)

	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	log.Printf("list of reges: %v\n\n\n\n\n\n", regs)

	results := new([]Registration)
	err = client.Call("Search", "email='dummy9001'", results)

	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	log.Printf("results of search: %v\n", results)

	asyncresults := new([]Registration)
	response := client.Async("Search", "email='dummy9001'", asyncresults)

	cb := <-response.Done

	if cb.Error != nil {
		log.Println("ERROR: " + cb.Error.Error())
	}

	log.Printf("results of async search: %v\n", asyncresults)

}
