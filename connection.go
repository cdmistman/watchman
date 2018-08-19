package watchman

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

type Request interface {
	Args() []interface{}
}

type Response interface {
	Version() string
	Warning() string
}

type response struct {
	Version string
	Warning string
}

type Connection struct {
	reader *bufio.Reader
	socket net.Conn
	// metadata
	capabilities map[string]struct{}
	sockname     string
	version      string
}

func Connect() (*Connection, error) {
	sockname, err := sockname()
	if err != nil {
		return nil, err
	}

	socket, err := dial(sockname, 30*time.Second)
	if err != nil {
		return nil, err
	}

	c := &Connection{
		reader:   bufio.NewReader(socket),
		socket:   socket,
		sockname: sockname,
	}
	err = c.init()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) HasCapability(capability string) bool {
	_, ok := c.capabilities[capability]
	return ok
}

func (c *Connection) SockName() string {
	return c.sockname
}

func (c *Connection) Version() string {
	return c.version
}

func (c *Connection) command(args ...interface{}) (resp map[string]interface{}, err error) {
	req, err := json.Marshal(args)
	if err != nil {
		return
	}

	_, err = fmt.Fprintln(c.socket, string(req))
	if err != nil {
		return
	}

	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return
	}

	err = json.Unmarshal(line, &resp)
	return
}

func (c *Connection) init() (err error) {
	resp, err := c.command("list-capabilities")
	if err != nil {
		return
	}

	if version, ok := resp["version"].(string); ok {
		c.version = version
	}

	if capabilities, ok := resp["capabilities"].([]interface{}); ok {
		capset := map[string]struct{}{}
		for _, cap := range capabilities {
			capset[cap.(string)] = struct{}{}
		}
		c.capabilities = capset
	}

	return
}

func (c *Connection) Recv(res interface{}) (Response, error) {
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var pdu map[string]json.RawMessage
	err = json.Unmarshal(line, &pdu)
	if msg, ok := pdu["error"]; ok {
		return nil, errors.New(string(msg)) // TODO custom Error
	} else if _, ok := pdu["subscription"]; ok {
		return nil, nil
	}

	err = json.Unmarshal(line, res)
	return nil, err
}

func (c *Connection) Send(req Request) (err error) {
	args := req.Args()
	b, err := json.Marshal(args)
	if err != nil {
		return
	}

	_, err = c.socket.Write(b)
	if err != nil {
		return
	}

	_, err = c.socket.Write([]byte("\n"))
	return
}

func sockname() (string, error) {
	sockname := os.Getenv("WATCHMAN_SOCK")
	if sockname != "" {
		return sockname, nil
	}

	buffer := &bytes.Buffer{}
	cmd := exec.Command("watchman", "get-sockname")
	cmd.Stdout = buffer
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var resp map[string]string
	if err := json.NewDecoder(buffer).Decode(&resp); err != nil {
		return "", err
	}

	sockname, ok := resp["sockname"]
	if !ok {
		return "", errors.New("unable to find watchman socket")
	}

	return sockname, nil
}
