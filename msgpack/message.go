package rpc

type Message struct {
	Key       string
	Namespace string
	Method    string
	Args      []interface{}
}
