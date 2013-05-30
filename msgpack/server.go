package rpc

import (
	"fmt"
	//"github.com/ugorji/go-msgpack"
	"github.com/ugorji/go/codec"
	"io"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"time"
)

type server struct {
	cn net.Conn
	ln net.Listener
}

var (
	mapStrIntfTyp = reflect.TypeOf(map[string]interface{}(nil))
	sliceByteTyp  = reflect.TypeOf([]byte(nil))
	timeTyp       = reflect.TypeOf(time.Time{})
)

// create and configure Handle
var (
	bh codec.BincHandle
	mh codec.MsgpackHandle
)

// create and use decoder/encoder
var (
	r io.Reader
	w io.Writer
	b []byte
	//h = &bh // or mh to use msgpack
	h = &mh
)

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
	svr.config()
	endpoint := fmt.Sprintf("%sRpc", svc.Name)
	rpc.RegisterName(endpoint, &svc)
	return svr.listen(port)
}
func (svr *server) config() {
	mh.MapType = mapStrIntfTyp

	// configure extensions for msgpack, to enable Binary and Time support for tags 0 and 1
	mh.AddExt(sliceByteTyp, 0, mh.BinaryEncodeExt, mh.BinaryDecodeExt)
	mh.AddExt(timeTyp, 1, mh.TimeEncodeExt, mh.TimeDecodeExt)

	//if wanting to enc/dec
	// dec := codec.NewDecoder(r, h)
	// dec = codec.NewDecoderBytes(b, h)
	// err := dec.Decode(&v)

	// enc := codec.NewEncoder(w, h)
	// enc = codec.NewEncoderBytes(&b, h)
	// err = enc.Encode(v)
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
	//rpcCodec := msgpack.NewCustomRPCServerCodec(conn, nil)
	rpcCodec := codec.MsgpackSpecRpc.ServerCodec(conn, h)
	go rpc.ServeCodec(rpcCodec)
}
func (svr *server) Stop() {
	defer svr.cn.Close()
}
