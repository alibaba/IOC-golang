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

type Mod struct {
	Prefix string
	Name   string
}

type Project struct {
	Path string
	Name string
	Mod  *Mod
}

func newMod() *Mod {
	return &Mod{}
}

func newProject() *Project {
	return &Project{
		Mod: newMod(),
	}
}

func initProject(opts ...Option) *Project {
	pt := newProject()
	for _, opt := range opts {
		opt(pt)
	}

	return pt
}

type Option func(*Project)

func WithPath(path string) Option {
	return func(project *Project) {
		project.Path = path
	}
}

func WithName(Name string) Option {
	return func(project *Project) {
		project.Name = Name
	}
}

func WithModPrefix(prefix string) Option {
	return func(project *Project) {
		project.Mod.Prefix = prefix
	}
}

func WithModName(name string) Option {
	return func(project *Project) {
		project.Mod.Name = name
	}
}
