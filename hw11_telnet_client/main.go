package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

var (
	timeout                     time.Duration
	ErrorMissingArgs            = errors.New("invalid args")
	ErrorUnsupportedHostAddress = errors.New("invalid host address")
	ErrorInvalidPort            = errors.New("invalid port number")
)

func init() {
	const defaultTimeout = 3 * time.Second
	pflag.DurationVarP(&timeout, "timeout", "t", defaultTimeout, "")
}

func getAddress(hostArg string, portArg string) (string, error) {
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

func main() {
	pflag.Parse()
	if pflag.NArg() < 2 {
		log.Fatal(ErrorMissingArgs)
	}

	address, err := getAddress(pflag.Arg(0), pflag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}

	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		log.Fatal(err)
	}

	defer tc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer stop()
		if err := tc.Send(); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	go func() {
		defer stop()
		if err := tc.Receive(); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	<-ctx.Done()
}
