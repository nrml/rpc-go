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
}

func NewServer(name string, svc interface{}, port int64) (server, error) {
	svr := server{}
	wrap, err := NewService(name, svc)
	if err != nil {
		log.Fatal("must provide a service name")
	}
	err = svr.init(&wrap, port)
	return svr, err
}

func (svr *server) init(svc *service, port int64) error {
	endpoint := fmt.Sprintf("%sRpc", svc.Name)
	rpc.RegisterName(endpoint, svc)
	svr.listen(port)
	return nil
}
func (svr *server) listen(port int64) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		conn, _ := l.Accept()

		svr.cn = conn
		log.Println("using custom codec for server")
		rpcCodec := msgpack.NewCustomRPCServerCodec(conn, nil)
		rpc.ServeCodec(rpcCodec)

	}
	return err
}
func (svr *server) Stop() {
	defer svr.cn.Close()
}
