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

package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alibaba/ioc-golang"
)

var Cmd = &cobra.Command{
	Use: "iocli",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("iocli version %s\n", ioc.Version)
			return
		}
		fmt.Println("hello")
	},
}

var versionFlag bool

func init() {
	Cmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Version of iocli")
}
