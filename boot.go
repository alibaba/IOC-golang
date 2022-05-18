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
	"github.com/fatih/color"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/debug"
)

func Load() error {
	printLogo()
	color.Cyan("Welcome to use ioc-golang!")

	// 1. load config
	color.Blue("[Boot] Start to load ioc-golang config")
	if err := config.Load(); err != nil {
		return err
	}

	// 2. load debug
	color.Blue("[Boot] Start to load debug")
	if err := debug.Load(); err != nil {
		return err
	}

	// 3. load autowire
	color.Blue("[Boot] Start to load autowire")
	return autowire.Load()
}
