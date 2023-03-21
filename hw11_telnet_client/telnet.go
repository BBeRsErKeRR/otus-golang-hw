package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrorNilConnection = errors.New("Unreachable")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (tc *telnetClient) transfer(in io.Reader, out io.Writer) error {
	if tc.conn == nil {
		return ErrorNilConnection
	}
	_, err := io.Copy(out, in)
	return err
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	return nil
}

func (tc *telnetClient) Send() error {
	return tc.transfer(tc.in, tc.conn)
}

func (tc *telnetClient) Receive() error {
	return tc.transfer(tc.conn, tc.out)
}

func (tc *telnetClient) Close() error {
	if tc.conn == nil {
		return nil
	}
	err := tc.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
