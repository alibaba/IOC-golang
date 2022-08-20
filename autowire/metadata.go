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

package autowire

// autowire metadata key

const MetadataKey = "autowire"

type AutowireMetadata map[string]interface{}

func ParseAutowireMetadataFromSDMetadata(metadata Metadata) AutowireMetadata {
	if metadata == nil {
		return nil
	}
	if autowireMetadataVal, ok := metadata[MetadataKey]; ok {
		if autowireMetadata, ok2 := autowireMetadataVal.(map[string]interface{}); ok2 {
			return autowireMetadata
		}
	}
	return nil
}

// common metadata keys

const CommonMetadataKey = "common"
const CommonImplementsMetadataKey = "implements"
const CommonActiveProfileMetadataKey = "activeProfile"

func parseCommonImplementsMetadataFromSDMetadata(metadata Metadata) []interface{} {
	autowireMetadata := ParseAutowireMetadataFromSDMetadata(metadata)
	if autowireMetadata == nil {
		return nil
	}
	autowireCommonMetadata, ok := autowireMetadata[CommonMetadataKey].(map[string]interface{})
	if !ok {
		return nil
	}
	result, ok := autowireCommonMetadata[CommonImplementsMetadataKey].([]interface{})
	if !ok {
		return nil
	}
	return result
}

func parseCommonActiveProfileMetadataFromSDMetadata(metadata Metadata) string {
	autowireMetadata := ParseAutowireMetadataFromSDMetadata(metadata)
	if autowireMetadata == nil {
		return ""
	}
	autowireCommonMetadata, ok := autowireMetadata[CommonMetadataKey].(map[string]interface{})
	if !ok {
		return ""
	}
	result, ok := autowireCommonMetadata[CommonActiveProfileMetadataKey].(string)
	if !ok {
		return ""
	}
	return result
}
