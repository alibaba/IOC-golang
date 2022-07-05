package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvPrefixKey = "${"
	EnvSuffixKey = "}"
	And          = "&"
)

func ExpandConfigEnvValue(targetValue interface{}) (interface{}, bool) {
	if tv, ok := targetValue.(string); ok {
		// ${ Xxx }
		if strings.HasPrefix(tv, EnvPrefixKey) && strings.HasSuffix(tv, EnvSuffixKey) && isEnv(tv) {
			expandValue := os.ExpandEnv(tv)
			if expandValue != "" {
				return expandValue, true
			}
		}
	}
	return targetValue, false
}

func ExpandConfigNestedValue(targetValue interface{}) (interface{}, bool) {
	if tv, ok := targetValue.(string); ok {
		if strings.HasPrefix(tv, EnvPrefixKey) && strings.HasSuffix(tv, EnvSuffixKey) && !isEnv(tv) {
			// try nested parsing
			var nestedValue interface{}
			// ${autowire.normal.<github.com/alibaba/ioc-golang/extension/state/redis.Redis>.expand.address}
			err := LoadConfigByPrefix(tv[2:len(tv)-1], &nestedValue)
			if err != nil {
				return nestedValue, false
			}

			return nestedValue, true
		}
	}
	return targetValue, false
}

func ExpandConfigValueIfNecessary(targetValue interface{}) interface{} {
	result, _ := ExpandConfigEnvValue(targetValue)
	result, _ = ExpandConfigNestedValue(result)
	return result
}

func parseEnvIfNecessary(config Config) {
	for k, v := range config {
		if val, ok := v.(string); ok {
			if expandValue, expand := ExpandConfigEnvValue(val); expand {
				config[k] = expandValue
				continue
			}
		} else if subConfig, ok := v.(Config); ok {
			parseEnvIfNecessary(subConfig)
		}
	}
}

func parseNestedIfNecessary(config Config) {
	for k, v := range config {
		if val, ok := v.(string); ok {
			if expandValue, expand := ExpandConfigNestedValue(val); expand {
				config[k] = expandValue
				continue
			}
		} else if subConfig, ok := v.(Config); ok {
			parseNestedIfNecessary(subConfig)
		}
	}
}

func parseConfigIfNecessary(config Config) {
	parseEnvIfNecessary(config)
	parseNestedIfNecessary(config)
}

func expandIfNecessary(targetValue string) string {
	// address=${REDIS_ADDRESS_EXPAND}&db=5
	if strings.Contains(targetValue, EnvPrefixKey) && strings.Contains(targetValue, EnvSuffixKey) {
		kvs := strings.Split(targetValue, And)
		kvz := make([]string, 0, len(kvs))
		for _, kv := range kvs {
			splitedKV := strings.Split(kv, "=")
			if len(splitedKV) != 2 {
				kvz = append(kvz, kv)
				continue
			}
			subKey := splitedKV[0]
			expandValue := ExpandConfigValueIfNecessary(splitedKV[1])
			kvz = append(kvz, fmt.Sprintf("%s=%s", subKey, expandValue))
		}

		targetValue = strings.Join(kvz, And)
	}

	return targetValue
}
