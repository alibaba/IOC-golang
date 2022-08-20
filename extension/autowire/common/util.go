package common

import "strings"

func ParseInterfacePkgAndInterfaceName(interfaceID string) (string, string) {
	splited := strings.Split(interfaceID, ".")
	return strings.Join(splited[:len(splited)-1], "."), splited[len(splited)-1]
}
