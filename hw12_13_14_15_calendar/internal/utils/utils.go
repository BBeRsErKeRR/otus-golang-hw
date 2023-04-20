package utils

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

var (
	ErrorUnsupportedProtocol    = errors.New("invalid protocol")
	ErrorUnsupportedHostAddress = errors.New("invalid host address")
	ErrorInvalidPort            = errors.New("invalid port number")
	addrMatcher                 = regexp.MustCompile(`^((([a-z0-9][a-z0-9\-]*[a-z0-9])|[a-z0-9])\.?)+$`)
)

func GetMqAddress(protocolArg, hostArg, portArg, userArg, passArg string) (string, error) {
	var address string

	if protocolArg != "amqp" {
		return address, ErrorUnsupportedProtocol
	}

	if hostArg != "localhost" && !addrMatcher.MatchString(hostArg) && net.ParseIP(hostArg) == nil {
		return address, ErrorUnsupportedHostAddress
	}

	port, err := strconv.Atoi(portArg)
	if err != nil {
		return address, err
	}

	if port < 1 || port > 65535 {
		return address, ErrorInvalidPort
	}

	netAddress := net.JoinHostPort(hostArg, portArg)

	if userArg != "" && passArg != "" {
		address = fmt.Sprintf("%v://%v:%v@%v", protocolArg, userArg, passArg, netAddress)
	} else {
		address = fmt.Sprintf("%v://%v", protocolArg, netAddress)
	}

	return address, nil
}

func GetAddress(hostArg string, portArg string) (string, error) {
	var address string

	if hostArg != "localhost" && !addrMatcher.MatchString(hostArg) && net.ParseIP(hostArg) == nil {
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
