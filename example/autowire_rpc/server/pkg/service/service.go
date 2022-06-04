package service

import (
	aliasDTO "github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/dto"
)

// +ioc:autowire=true
// +ioc:autowire:type=rpc

type ServiceStruct struct {
}

func (s *ServiceStruct) GetUser(name string, age int, age32 int32, age64 int64, ageF32 float32, ageF64 float64, agePtr *int, age32Ptr *int32, age64Ptr *int64, ageF32Ptr *float32, ageF64Ptr *float64, reqUsr *aliasDTO.User) (*aliasDTO.User, *aliasDTO.User, error) {
	return &aliasDTO.User{
		Name: name,
		Age:  age,
	}, reqUsr, nil
}
