/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package protocol_impl

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"dubbo.apache.org/dubbo-go/v3/common"
	perrors "github.com/pkg/errors"
)

// ServiceMap store description of service.
var ServiceMap = &serviceMap{
	serviceMap:   make(map[string]map[string]*Service),
	interfaceMap: make(map[string][]*Service),
}

type serviceMap struct {
	mutex        sync.RWMutex                   // protects the serviceMap
	serviceMap   map[string]map[string]*Service // protocol -> service name -> service
	interfaceMap map[string][]*Service          // interface -> service
}

// GetService gets a service definition by protocol and name
func (sm *serviceMap) GetService(protocol, interfaceName, group, version string) *Service {
	serviceKey := common.ServiceKey(interfaceName, group, version)
	return sm.GetServiceByServiceKey(protocol, serviceKey)
}

// GetService gets a service definition by protocol and service key
func (sm *serviceMap) GetServiceByServiceKey(protocol, serviceKey string) *Service {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	if s, ok := sm.serviceMap[protocol]; ok {
		if srv, ok := s[serviceKey]; ok {
			return srv
		}
		return nil
	}
	return nil
}

// GetInterface gets an interface definition by interface name
func (sm *serviceMap) GetInterface(interfaceName string) []*Service {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	if s, ok := sm.interfaceMap[interfaceName]; ok {
		return s
	}
	return nil
}

// Register registers a service by @interfaceName and @protocol
func (sm *serviceMap) Register(interfaceName, protocol, group, version string, rcvr common.RPCService) (string, error) {
	if sm.serviceMap[protocol] == nil {
		sm.serviceMap[protocol] = make(map[string]*Service)
	}
	if sm.interfaceMap[interfaceName] == nil {
		sm.interfaceMap[interfaceName] = make([]*Service, 0, 16)
	}

	s := new(Service)
	s.rcvrType = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)

	sname := common.ServiceKey(interfaceName, group, version)
	if server := sm.GetService(protocol, interfaceName, group, version); server != nil {
		return "", perrors.New("service already defined: " + sname)
	}
	s.methods = make(map[string]*MethodType)

	// Install the methods
	methods := ""
	methods, s.methods = suitableMethods(s.rcvrType)

	sm.mutex.Lock()
	sm.serviceMap[protocol][sname] = s
	sm.interfaceMap[interfaceName] = append(sm.interfaceMap[interfaceName], s)
	sm.mutex.Unlock()

	return strings.TrimSuffix(methods, ","), nil
}

// Service is description of service
type Service struct {
	name     string
	rcvr     reflect.Value
	rcvrType reflect.Type
	methods  map[string]*MethodType
}

// Method gets @s.methods.
func (s *Service) Method() map[string]*MethodType {
	return s.methods
}

// Name will return service name
func (s *Service) Name() string {
	return s.name
}

// RcvrType gets @s.rcvrType.
func (s *Service) RcvrType() reflect.Type {
	return s.rcvrType
}

// Rcvr gets @s.rcvr.
func (s *Service) Rcvr() reflect.Value {
	return s.rcvr
}

// suitableMethods returns suitable Rpc methods of typ
func suitableMethods(typ reflect.Type) (string, map[string]*MethodType) {
	methods := make(map[string]*MethodType)
	var mts []string
	method, ok := typ.MethodByName(common.METHOD_MAPPER)
	var methodMapper map[string]string
	if ok && method.Type.NumIn() == 1 && method.Type.NumOut() == 1 && method.Type.Out(0).String() == "map[string]string" {
		methodMapper = method.Func.Call([]reflect.Value{reflect.New(typ.Elem())})[0].Interface().(map[string]string)
	}

	for m := 0; m < typ.NumMethod(); m++ {
		method = typ.Method(m)
		if mt := suiteMethod(method); mt != nil {
			methodName, ok := methodMapper[method.Name]
			if !ok {
				methodName = method.Name
			}
			methods[methodName] = mt
			mts = append(mts, methodName)
		}
	}
	return strings.Join(mts, ","), methods
}

// suiteMethod returns a suitable Rpc methodType
func suiteMethod(method reflect.Method) *MethodType {
	mtype := method.Type
	inNum := mtype.NumIn()
	outNum := mtype.NumOut()

	// Method must be exported.
	if method.PkgPath != "" {
		return nil
	}

	var (
		replyType, ctxType reflect.Type
		argsType           []reflect.Type
	)

	// replyType
	if outNum == 2 {
		replyType = mtype.Out(0)
		if !isExportedOrBuiltinType(replyType) {
			return nil
		}
	}

	index := 1

	// ctxType
	if inNum > 1 && mtype.In(1).String() == "context.Context" {
		ctxType = mtype.In(1)
		index = 2
	}

	for ; index < inNum; index++ {
		argsType = append(argsType, mtype.In(index))
		// need not be a pointer.
		if !isExportedOrBuiltinType(mtype.In(index)) {
			return nil
		}
	}

	return &MethodType{method: method, argsType: argsType, replyType: replyType, ctxType: ctxType}
}

// MethodType is description of service method.
type MethodType struct {
	method    reflect.Method
	ctxType   reflect.Type   // request context
	argsType  []reflect.Type // args except ctx, include replyType if existing
	replyType reflect.Type   // return value, otherwise it is nil
}

// Method gets @m.method.
func (m *MethodType) Method() reflect.Method {
	return m.method
}

// CtxType gets @m.ctxType.
func (m *MethodType) CtxType() reflect.Type {
	return m.ctxType
}

// ArgsType gets @m.argsType.
func (m *MethodType) ArgsType() []reflect.Type {
	return m.argsType
}

// ReplyType gets @m.replyType.
func (m *MethodType) ReplyType() reflect.Type {
	return m.replyType
}

// SuiteContext transfers @ctx to reflect.Value type or get it from @m.ctxType.
func (m *MethodType) SuiteContext(ctx context.Context) reflect.Value {
	if ctxV := reflect.ValueOf(ctx); ctxV.IsValid() {
		return ctxV
	}
	return reflect.Zero(m.ctxType)
}

// Is this an exported - upper case - name
func isExported(name string) bool {
	s, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(s)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
