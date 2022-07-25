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

package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/autowire"
)

func TestParseRollbackMethodNameFromSDMetadata(t *testing.T) {
	t.Run("parse success", func(t *testing.T) {
		metadata := make(autowire.Metadata)
		metadata["aop"] = map[string]interface{}{
			"transaction": map[string]string{
				"methodName":  "rollbackName",
				"methodName2": "",
			},
		}
		val, ok := parseRollbackMethodNameFromSDMetadata(metadata, "methodName")
		assert.True(t, ok)
		assert.Equal(t, "rollbackName", val)

		val, ok = parseRollbackMethodNameFromSDMetadata(metadata, "methodName2")
		assert.True(t, ok)
		assert.Equal(t, "", val)
	})

	t.Run("parse with tx metadata invalid type", func(t *testing.T) {
		metadata := make(autowire.Metadata)
		metadata["aop"] = map[string]interface{}{
			"transaction": map[string]interface{}{
				"methodName":  "rollbackName",
				"methodName2": "",
			},
		}
		val, ok := parseRollbackMethodNameFromSDMetadata(metadata, "methodName")
		assert.True(t, !ok)
		assert.Equal(t, "", val)
	})

	t.Run("parse with aop metadata empty", func(t *testing.T) {
		metadata := make(autowire.Metadata)
		metadata["autowire"] = map[string]interface{}{
			"transaction": map[string]string{
				"methodName":  "rollbackName",
				"methodName2": "",
			},
		}
		val, ok := parseRollbackMethodNameFromSDMetadata(metadata, "methodName")
		assert.True(t, !ok)
		assert.Equal(t, "", val)
	})

}
