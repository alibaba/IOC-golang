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

package ioc

import (
	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/logger"

	_ "github.com/alibaba/ioc-golang/extension/imports/boot"
)

const Version = "1.0.4"

func Load(opts ...config.Option) error {
	printLogo()
	logger.Cyan("Welcome to use ioc-golang %s!", Version)

	// 1. load config
	logger.Blue("[Boot] Start to load ioc-golang config")
	if err := config.Load(opts...); err != nil {
		return err
	}

	// 2. load debug
	logger.Blue("[Boot] Start to load debug")
	if err := aop.Load(); err != nil {
		return err
	}

	// 3. load autowire
	logger.Blue("[Boot] Start to load autowire")
	return autowire.Load()
}
