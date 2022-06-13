package common

import (
	"net"
	"strconv"

	"github.com/fatih/color"
)

func GetTCPListener(port string) (net.Listener, error) {
	lst, err := net.Listen("tcp", ":"+port)
	for err != nil {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			color.Blue("[Debug] Debug server listening with invalid port :%s, error = %s", port, err)
			return nil, err
		}
		color.Blue("[Debug] Debug server listening port :%s failed with error = %s, try to bind %d", port, err, portInt+1)
		port = strconv.Itoa(portInt + 1)
		lst, err = net.Listen("tcp", ":"+port)
	}
	return lst, nil
}
