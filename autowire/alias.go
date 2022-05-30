package autowire

import (
	"fmt"
)

// sdIDAliasMap a map of SDID alias-mapping named as alias container.
var sdIDAliasMap = make(map[string]string)

func RegisterAlias(alias, value string) {
	if _, ok := sdIDAliasMap[alias]; ok {
		panic(fmt.Sprintf("[Autowire Alias] Duplicate alias:[%s]", alias))
	}

	sdIDAliasMap[alias] = value
}

func MappingSDIDAliasIfNecessary(sdID string) string {
	if mappingSDID, ok := getAlias(sdID); ok {
		return mappingSDID
	}

	return sdID
}

func HasAlias(alias string) bool {
	_, ok := sdIDAliasMap[alias]

	return ok
}

func getAlias(alias string) (string, bool) {
	v, ok := sdIDAliasMap[alias]

	return v, ok
}
