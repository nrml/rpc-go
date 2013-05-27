package rpc

import (
	"fmt"
	"log"
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

	log.Printf("list of reges: %v\n", regs)

}
