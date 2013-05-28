package rpc

import (
	"fmt"
	"github.com/ugorji/go-msgpack"
	"log"
	"net"
	"net/rpc"
)

type server struct {
	cn net.Conn
	ln net.Listener
}

func NewServer(name string, svc interface{}, port int64) (server, error) {
	svr := server{}
	wrap, err := NewService(name, svc)
	if err != nil {
		log.Fatal("must provide a service name")
	}
	err = svr.init(wrap, port)
	return svr, err
}

func (svr *server) init(svc service, port int64) error {
	endpoint := fmt.Sprintf("%sRpc", svc.Name)
	rpc.RegisterName(endpoint, &svc)
	return svr.listen(port)
}
func (svr *server) listen(port int64) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatal("listen error:", err)
	}
	svr.ln = l

	return err
}
func (svr *server) Accept() {
	conn, _ := svr.ln.Accept()

	svr.cn = conn
	log.Println("using custom codec for server")
	rpcCodec := msgpack.NewCustomRPCServerCodec(conn, nil)
	go rpc.ServeCodec(rpcCodec)
}
func (svr *server) Stop() {
	defer svr.cn.Close()
}
