package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
)

const (
	vmmBinary = "qemu-system-x86_64"
	vmmSocket = "/tmp/qmp.sock"
	unixNtw   = "unix"
)

var (
	defArgs = []string{
		"-M", "q35",
		"-m", "512M",
		"-display", "none",
		"-nodefaults",
		"-qmp", "unix:/tmp/qmp.sock,server,wait=off",
	}
)

type Qemu struct {
	conn   net.Conn
	reader *bufio.Reader
}

func StartServer() (*exec.Cmd, error) {

	cmd := exec.Command(vmmBinary, defArgs...)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}


func (q *Qemu) Close() error {
	
	return q.conn.Close()
}

func Connect() (*Qemu, error) {
	conn, err := net.Dial(unixNtw, vmmSocket)
	if err != nil {
		return nil, err
	}

	q := &Qemu{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}

	return q, nil
}

func (q *Qemu) Read(ctx context.Context, done chan interface{})  {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Read goroutine closed.")
			close(done)
			return 

		default:
			msg, err := q.reader.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(msg))
		}
	}
}

func (q *Qemu) Execute(cmd map[string]interface{}) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = q.conn.Write(data)

	return err
}
