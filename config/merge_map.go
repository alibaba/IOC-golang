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

package config

import (
	"reflect"

	"github.com/fatih/color"
)

const (
	defaultMergeDepth uint8 = 1 << 3
	maxMergeDepth     uint8 = 1 << 4
)

type AnyMap = map[string]interface{} // alias

// MergeMap
//
// @dangerous trigger PANIC
//
// @return dst
func MergeMap(dst, src AnyMap, maxDepths ...uint8) AnyMap {
	return merge(dst, src, 0, maxDepths...)
}

func merge(dst, src AnyMap, depth uint8, depths ...uint8) AnyMap {
	maxDepth := determineMerDepth(depths)
	if maxDepth > maxMergeDepth {
		panic(color.RedString("[Config] expect depth too deep: [%d]", maxDepth))
	}
	color.Blue("[Config] merge config map, depth: [%d]", depth)
	if depth > maxDepth {
		panic(color.RedString("[Config] recursion too deep: [%d]", depth))
	}

	for k, v := range src {
		if dv, ok := dst[k]; ok {
			dstMap, dstOk := toMap(dv)
			srcMap, srcOk := toMap(v)
			if srcOk && dstOk {
				v = merge(dstMap, srcMap, depth+1, depths...)
			}
		}

		dst[k] = v
	}

	return dst
}

func determineMerDepth(depths []uint8) uint8 {
	depth := defaultMergeDepth
	switch len(depths) {
	case 1:
		depth = depths[0]
	}

	return depth
}

func toMap(src interface{}) (AnyMap, bool) {
	value := reflect.ValueOf(src)
	if value.Kind() == reflect.Map {
		am := AnyMap{}
		for _, k := range value.MapKeys() {
			am[k.String()] = value.MapIndex(k).Interface()
		}
		return am, true
	}

	return AnyMap{}, false
}
