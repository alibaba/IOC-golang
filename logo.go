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
)

var logo = "  ___    ___     ____                           _                         \n |_ _|  / _ \\   / ___|           __ _    ___   | |   __ _   _ __     __ _ \n  | |  | | | | | |      _____   / _` |  / _ \\  | |  / _` | | '_ \\   / _` |\n  | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |\n |___|  \\___/   \\____|          \\__, |  \\___/  |_|  \\__,_| |_| |_|  \\__, |\n                                |___/                               |___/ "

func printLogo() {
	color.Cyan(logo)
}
