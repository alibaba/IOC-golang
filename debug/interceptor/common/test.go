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

package common

import (
	"context"
)

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
