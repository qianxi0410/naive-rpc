package router

import "context"

type HandldFunc = func(svr interface{}, ctx context.Context, req interface{}) (rsp interface{}, err error)

// method description
type MethodDesc struct {
	Name   string
	Method HandldFunc
}

// stream scription
type StreamDesc struct {
}

// service description
type ServiceDesc struct {
	// like serviceName.serviceMethod
	Name   string
	Type   interface{}
	Method map[string]*MethodDesc
	Stream map[string]*StreamDesc
}
