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

func GetSDIDByAliasIfNecessary(key string) string {
	if mappingSDID, ok := GetSDIDByAlias(key); ok {
		return mappingSDID
	}

	return key
}

func GetSDIDByAlias(alias string) (string, bool) {
	sdid, ok := sdIDAliasMap[alias]
	return sdid, ok
}
