package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type TClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (t *TClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	t.conn = conn
	return nil
}

func (t *TClient) Close() error {
	return t.conn.Close()
}

func (t *TClient) Send() error {
	var err error
	for sendErrorsCount := 0; sendErrorsCount < 3; sendErrorsCount++ {
		_, err = io.Copy(t.conn, t.in)
		if err == nil {
			break
		}
	}
	return err
}

func (t *TClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}
