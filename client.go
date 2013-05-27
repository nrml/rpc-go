package rpc

import (
	"errors"
	"fmt"
	"github.com/ugorji/go-msgpack"
	"log"
	"net"
	"net/rpc"
)

type client struct {
	key       string
	namespace string
	service   string
	rpc       *rpc.Client
}

func NewClient(key string, namespace string, svcname string) (client, error) {
	var err error
	c := client{}
	if key == "" || namespace == "" {
		err = errors.New("needs secret and namespace")
	}
	c.service = svcname
	c.key = key
	c.namespace = namespace
	return c, err
}

func (c *client) Connect(host string, port int64) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	log.Println("using custom codec for client")
	rpcCodec := msgpack.NewCustomRPCClientCodec(conn, nil)
	client := rpc.NewClientWithCodec(rpcCodec)
	c.rpc = client
	return err
}

//last agrument is the reply that is a zeroed T value pointer (new), all others sent with the request
func (c *client) Call(method string, args ...interface{}) error {
	ep := c.endpoint()
	l := len(args)
	reqargs := args[0 : l-1]
	reply := args[l-1:][0]
	msg := Message{c.key, c.namespace, method, reqargs}
	var err error
	err = c.rpc.Call(ep, msg, reply)
	return err
}

func (c *client) Async(method string, args ...interface{}) *rpc.Call {
	ep := c.endpoint()
	l := len(args)
	reqargs := args[0 : l-1]
	reply := args[l-1:][0]
	msg := Message{c.key, c.namespace, method, reqargs}

	return c.rpc.Go(ep, msg, reply, nil)
}
func (c *client) endpoint() string {
	return fmt.Sprintf("%sRpc.Call", c.service)
}
