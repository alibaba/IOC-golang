package common

import "context"

type User struct {
	Name string
}

type RequestParam struct {
	User *User
}

type Response struct {
	Name string
}

type ServiceFoo struct {
}

func (s *ServiceFoo) Invoke(ctx context.Context, param *RequestParam) (*Response, error) {
	return &Response{
		Name: param.User.Name,
	}, nil
}
