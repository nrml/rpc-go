package rpc

import (
	"errors"
	"github.com/nrml/convert-go"
	"log"
	"reflect"
)

type service struct {
	Name      string
	Object    interface{}
	Key       string
	Namespace string
}

func NewService(name string, local interface{}) (service, error) {
	var err error
	if name == "" {
		err = errors.New("requires a service name")
	}
	svc := service{name, local, "", ""}
	return svc, err
}

func (svc *service) Call(msg Message, reply *interface{}) error {
	//log.Println("made it into call")

	if msg.Key == "" || msg.Namespace == "" {
		return errors.New("must have key and namespace.")
	}

	var err error
	if msg.Key == "" || msg.Namespace == "" {
		err = errors.New("key and namespace required")
		return err
	}

	//log.Println("about to relect on svc.Object")

	ro := reflect.ValueOf(svc.Object)

	//log.Println("about to reflect on key")
	sf := ro.Elem().FieldByName("Key")
	if sf.CanSet() {
		sf.SetString(msg.Key)
	}
	//log.Println("about to reflect on namespace")
	nsf := ro.Elem().FieldByName("Namespace")
	if nsf.CanSet() {
		nsf.SetString(msg.Namespace)
	}

	//initialize service
	minit := ro.MethodByName("Init")
	sandn := reflect.ValueOf(msg.Key + "." + msg.Namespace)
	minit.Call([]reflect.Value{sandn})

	//log.Println("about to try and get method")
	m := ro.MethodByName(msg.Method)
	mtype := m.Type()
	//log.Printf("about to to call the function(%v) with original args: %v\n", msg.Method, msg.Args)

	args := make([]reflect.Value, len(msg.Args))

	//TODO: there should be a better way to do this
	for i := 0; i < len(msg.Args); i++ {
		arg := msg.Args[i]
		argtype := mtype.In(i)
		log.Printf("type: %v", argtype)
		cp := reflect.New(argtype)
		k := reflect.ValueOf(arg).Kind().String()
		log.Printf("kind: %v\n", k)
		switch k {
		case "map":
			log.Println("convert map")
			mp := arg.(map[interface{}]interface{})
			convert.ConvertMap(mp, cp.Interface())
		case "struct":
			log.Println("convert struct")
			convert.Convert(arg, cp.Interface())
		default:
			if reflect.TypeOf(arg).ConvertibleTo(argtype) {
				log.Println("convert default: %v=%v", cp, arg)
				resp := reflect.ValueOf(arg).Convert(argtype)
				cp.Elem().Set(resp)
			}
		}

		args[i] = cp.Elem()
	}

	//log.Printf("about to to call the function(%v) with converted args: %v\n", msg.Method, len(args))
	resp := m.Call(args)

	log.Printf("made call, got response: %v", resp)
	//response is idx0, error idx1
	robj := resp[0].Interface()
	log.Printf("sending back %v\n", robj)
	ierr := resp[1].Interface()

	if ierr != nil {
		err = ierr.(error)
	}

	//log.Printf("first response obj: %v\n", robj)

	*reply = robj

	return err

}
