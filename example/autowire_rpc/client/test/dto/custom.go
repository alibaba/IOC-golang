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

package dto

type CustomStruct struct {
	User                  // combination
	CustomStructId        int64
	IdPtr                 *int64
	CustomStructName      string
	NamePtr               *string
	CustomStringMap       map[string]string
	CustomIntMap          map[string]int
	CustomSubStructPtrMap map[string]*User
	CustomSubStructMap    map[string]User
	StringSlice           []string
	SubStructSlice        []User
	SubStructPtrSlice     []*User
	SubStruct             User
	SubStructPtr          *User
}

func (c *CustomStruct) GetUser() User {
	return c.SubStruct
}

type UserGetter interface {
	GetUser() User
}
