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
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/config"
)

type Config struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	TableName string
}

func (c *Config) New(mysqlImpl *Impl) (*Impl, error) {
	var err error
	mysqlImpl.db, err = gorm.Open(mysql.Open(getMysqlLinkStr(c)), &gorm.Config{})
	mysqlImpl.tableName = c.TableName
	return mysqlImpl, err
}

func getMysqlLinkStr(conf *Config) string {
	return conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.DBName +
		"?charset=utf8&parseTime=True&loc=Local"
}

type paramLoader struct {
}

func (p *paramLoader) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	splitedTagValue := strings.Split(fi.TagValue, ",")
	param := &Config{}
	if len(splitedTagValue) <= 2 {
		return nil, fmt.Errorf("file info %s doesn't contains param information, create param from sd paramLoader failed", fi)
	}
	if err := config.LoadConfigByPrefix(fmt.Sprintf("autowire%[1]snormal%[1]sgithub.com/alibaba/ioc-golang/extension/normal/mysql.Impl%[1]s%[2]s%[1]sparam", config.YamlConfigSeparator, splitedTagValue[1]), param); err != nil {
		return nil, err
	}
	param.TableName = splitedTagValue[2]
	return param, nil
}
