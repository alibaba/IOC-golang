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

package mysql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getMysqlLinkStr(t *testing.T) {
	conf := &Config{
		Host:     "192.168.1.1",
		Port:     "1234",
		Username: "admin",
		Password: "admin",
		DBName:   "mydb",
	}
	got := getMysqlLinkStr(conf)
	assert.Equal(t, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName), got)
}
