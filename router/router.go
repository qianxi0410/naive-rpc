package router

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/kpango/glg"
	"github.com/qianxi0410/naive-rpc/errors"
)

type HandleWrapper func(ctx context.Context, req interface{}) (rsp interface{}, err error)

type Router struct {
	mapping map[string]HandleWrapper
	mutex   *sync.RWMutex
}

// return a router
func NewRouter() *Router {
	return &Router{
		mapping: make(map[string]HandleWrapper),
		mutex:   &sync.RWMutex{},
	}
}

// register service
func (r *Router) RegisterService(srvDesc *ServiceDesc, srvImpl interface{}) error {
	t := reflect.TypeOf(srvDesc.Type).Elem()
	i := reflect.TypeOf(srvImpl)

	if !i.Implements(t) {
		return fmt.Errorf("%s not implements interface %s", i, t)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, m := range srvDesc.Method {
		// gretting.Helloworld
		rpc := srvDesc.Name + "/" + m.Name
		f := func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
			// handle func
			return m.Method(srvImpl, ctx, req)
		}
		r.mapping[rpc] = f
	}

	return nil
}

func (r *Router) Forward(rpcName string, handleFunc HandleWrapper) {
	glg.Successf("%s router is registed", rpcName)
	r.mapping[rpcName] = handleFunc
}

func (r *Router) Route(rpcName string) (HandleWrapper, error) {
	f, ok := r.mapping[rpcName]
	if !ok {
		return nil, errors.RouteNotFoundErr
	}
	return f, nil
}
