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

package init

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr/v2"
)

//
// Nodes:
//
// The implementation of `init` cmd refers to Jupiter[Jupiter Github](https://github.com/douyu/jupiter)(v0.4.0) `new` cmd,
// but some enhancements have been made.
// For example: support for custom parameters and environment variables.
//

var project *Project

func Run(opts ...Option) {
	project = initProject(opts...)
	err := createProject()
	if err != nil {
		fmt.Printf("init failed: %v", err)
	}
}

func createProject() error {
	if err := handleProperties(); err != nil {
		return err
	}

	if err := doCreate(); err != nil {
		return err
	}

	doReport()

	return nil
}

func handleProperties() (err error) {
	if project.Path != "" {
		project.Path = parseEVNecessary(filepath.Clean(project.Path))

		if project.Path, err = filepath.Abs(project.Path); err != nil {
			return
		}
		project.Path = filepath.Join(project.Path, project.Name)
	} else {
		pwd, _ := os.Getwd()
		project.Path = filepath.Join(pwd, project.Name)
	}

	resetModIfNecessary()

	return
}

func resetModIfNecessary() {
	if project.Mod.Prefix == "" || strings.TrimSpace(project.Mod.Prefix) == "" {
		modPath := determineModPath(project.Path)
		project.Mod.Prefix = modPath
	}
	if project.Mod.Name == "" || strings.TrimSpace(project.Mod.Name) == "" {
		project.Mod.Name = fmt.Sprintf("%s/%s", project.Mod.Prefix, project.Name)
	}
}

func doReport() {
	fmt.Println("---------------- iocli init report ----------------")
	fmt.Println("Project dir:", project.Path)
	fmt.Println("Run:")
	fmt.Println("$ cd " + project.Path)
	fmt.Println("$ go mod tidy [-go=1.16 && go mod tidy -go=1.17]")
	fmt.Println("---------------- iocli init report ----------------")
}

//go:generate packr2
func doCreate() (err error) {
	box := packr.New("all", "./templates")
	if err = os.MkdirAll(project.Path, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		tmpl, _ := box.FindString(name)
		i := strings.LastIndex(name, string(os.PathSeparator))
		if i > 0 {
			dir := name[:i]
			if err = os.MkdirAll(filepath.Join(project.Path, dir), 0755); err != nil {
				return
			}
		}
		name = strings.TrimSuffix(name, ".tmpl")
		if err = doWriteFile(filepath.Join(project.Path, name), tmpl); err != nil {
			return
		}
	}

	return
}

func doWriteFile(path, tmpl string) (err error) {
	data, err := parseTmpl(tmpl)
	if err != nil {
		return
	}
	fmt.Println("File -> ", path)

	return os.WriteFile(path, data, 0755)
}

func parseTmpl(tmpl string) ([]byte, error) {
	tmp, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = tmp.Execute(&buf, project); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func determineModPath(projectPath string) (modPath string) {
	dir := filepath.Dir(projectPath)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			content, _ := os.ReadFile(filepath.Join(dir, "go.mod"))
			mod := find(`module\s+(?P<name>[\S]+)`, string(content), "$name")
			name := strings.TrimPrefix(filepath.Dir(projectPath), dir)
			name = strings.TrimPrefix(name, string(os.PathSeparator))
			if name == "" {
				return fmt.Sprintf("%s/", mod)
			}

			return fmt.Sprintf("%s/%s/", mod, name)
		}
		parent := filepath.Dir(dir)
		if dir == parent {
			return ""
		}

		dir = parent
	}
}

func find(regex, src, name string) string {
	var result []byte
	pattern := regexp.MustCompile(regex)
	for _, sub := range pattern.FindAllStringSubmatchIndex(src, -1) {
		result = pattern.ExpandString(result, name, src, sub)
	}

	return string(result)
}

func parseEVNecessary(path string) string {
	dirs := strings.Split(path, string(os.PathSeparator))
	nodes := make([]string, 0, len(dirs))
	for _, node := range dirs {
		if strings.HasPrefix(node, "$") {
			env := os.Getenv(node[1:])
			if env != "" {
				nodes = append(nodes, env)
				continue
			}
		}
		nodes = append(nodes, node)
	}

	return strings.Join(nodes, string(os.PathSeparator))
}
