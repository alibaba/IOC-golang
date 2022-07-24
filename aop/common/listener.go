package common

import (
	"net"
	"strconv"

	"github.com/alibaba/ioc-golang/logger"
)

func GetTCPListener(port string) (net.Listener, error) {
	lst, err := net.Listen("tcp", ":"+port)
	for err != nil {
		portInt, iToAError := strconv.Atoi(port)
		if iToAError != nil {
			logger.Blue("[Debug] Debug server listening with invalid port :%s, error = %s", port, iToAError)
			return nil, iToAError
		}
		logger.Blue("[Debug] Debug server listening port :%s failed with error = %s, try to bind %d", port, err, portInt+1)
		port = strconv.Itoa(portInt + 1)
		lst, err = net.Listen("tcp", ":"+port)
	}
	return lst, nil
}
