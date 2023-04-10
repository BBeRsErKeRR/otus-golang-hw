package server

import (
	"errors"
	"net"
	"regexp"
	"strconv"
)

var (
	ErrorUnsupportedHostAddress = errors.New("invalid host address")
	ErrorInvalidPort            = errors.New("invalid port number")
)

func GetAddress(hostArg string, portArg string) (string, error) {
	var address string
	re := regexp.MustCompile(`^((([a-z0-9][a-z0-9\-]*[a-z0-9])|[a-z0-9])\.?)+$`)

	if hostArg != "localhost" && !re.MatchString(hostArg) && net.ParseIP(hostArg) == nil {
		return address, ErrorUnsupportedHostAddress
	}

	port, err := strconv.Atoi(portArg)
	if err != nil {
		return address, err
	}

	if port < 1 || port > 65535 {
		return address, ErrorInvalidPort
	}

	address = net.JoinHostPort(hostArg, portArg)
	return address, nil
}
