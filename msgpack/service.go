package rpc

import (
	"errors"
	"fmt"
	"github.com/nrml/convert-go"
	"log"
	"reflect"
	"time"
)

type service struct {
	Name       string
	object     interface{}
	servicemap map[string]interface{}
	timeoutmap map[string]*time.Timer
}

func NewService(name string, local interface{}) (service, error) {
	log.Println("new service")
	var err error
	if name == "" {
		err = errors.New("requires a service name")
	}
	svcmap := make(map[string]interface{})
	tomap := make(map[string]*time.Timer)
	svc := service{name, local, svcmap, tomap}
	return svc, err
}

func (svc *service) Call(msg Message, reply *interface{}) error {
	//log.Println("call")
	if msg.Key == "" || msg.Namespace == "" {
		return errors.New("must have key and namespace.")
	}

	var err error
	if msg.Key == "" || msg.Namespace == "" {
		err = errors.New("key and namespace required")
		return err
	}

	var ro reflect.Value

	//find existing service in map
	chk := fmt.Sprintf("%s.%s", msg.Key, msg.Namespace)
	exst, ok := svc.servicemap[chk]

	if !ok {
		ot := reflect.ValueOf(svc.object).Elem().Type()
		exst = reflect.New(ot).Interface()
		//add service to map
		svc.servicemap[chk] = exst
	}

	ro = reflect.ValueOf(exst)

	m := ro.MethodByName(msg.Method)

	if !m.IsValid() {
		err = errors.New("unknown service method")
		return err
	}

	mtype := m.Type()

	l := len(msg.Args)
	li := mtype.NumIn()

	if l != li {
		err = errors.New(fmt.Sprintf("arg number mismatch: send %d, need %d\n", l, li))
	}

	args := make([]reflect.Value, l)

	//TODO: there should be a better way to do this
	for i := 0; i < len(msg.Args); i++ {
		arg := msg.Args[i]
		desttype := mtype.In(i)

		cp := reflect.New(desttype)
		k := reflect.ValueOf(arg).Kind().String()

		switch k {
		case "map":
			mp := arg.(map[string]interface{})
			convert.ConvertMap(mp, cp.Interface())
		case "struct":
			convert.Convert(arg, cp.Interface())
		default:
			if reflect.TypeOf(arg).ConvertibleTo(desttype) {
				resp := reflect.ValueOf(arg).Convert(desttype)
				cp.Elem().Set(resp)
			} else {
				log.Printf("cannot convert %v to function argument type %v\n", arg, desttype)
			}
		}

		args[i] = cp.Elem()
	}

	resp := m.Call(args)

	robj := resp[0].Interface()
	ierr := resp[1].Interface()

	if ierr != nil {
		err = ierr.(error)
	}

	*reply = robj

	timer, ok := svc.timeoutmap[chk]
	if ok {
		log.Println("resetting timer")
		timer.Reset(100 * time.Millisecond)
	} else {
		log.Println("new timer")
		timer = time.NewTimer(100 * time.Millisecond)
		svc.timeoutmap[chk] = timer
	}

	go svc.remove(timer, chk)

	return err
}
func (svc *service) remove(t *time.Timer, ns string) {
	_ = <-t.C
	t.Stop()
	log.Println("deleting service")
	delete(svc.timeoutmap, ns)
}
