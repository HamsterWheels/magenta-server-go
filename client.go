package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	nickname    string
	realname    string
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	serverInput chan Message
	output      chan string
	idle        time.Time
}

// NewClient returns a pointer to a Client object
func NewClient(name string, conn net.Conn, serverInput chan Message) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	client := &Client{
		nickname:    name,
		realname:    "",
		conn:        conn,
		reader:      reader,
		writer:      writer,
		serverInput: serverInput,
		output:      make(chan string),
		idle:        time.Now(),
	}

	// Start the loops that constantly monitor for input and output
	// from the client
	client.run()

	return client
}

// run starts two loops, concurrently, that monitor the Client's input and
// output channels
func (c *Client) run() {
	go c.processInput()
	go c.SendOutput()
}

// GetInput receives data from the Client and sends it to the server
func (c *Client) processInput() {
	for {
		// The buffer will stop reading when it detects a new line chatacter
		input, _ := c.reader.ReadString('\n')
		c.idle = time.Now()
		// If the client sends an empty message just ignore it
		if input := trimMessage(input); !isEmpty(input) {
			msg := NewMessage(c, input)
			c.serverInput <- *msg

		}
	}
}

// Close closes the clients connection and frees resources
func (c *Client) Close(msg string) {
	c.output <- msg
	c.conn.Close()
	close(c.output)
}

// SendOutput monitors the output channel.  When data is received from the server
// it is sent to the Client
func (c *Client) SendOutput() {
	for data := range c.output {
		c.writer.WriteString(data)
		c.writer.Flush()
	}
}

// isIdle checks the last time the Client was active.
func (c *Client) isIdle() bool {
	now := time.Now()
	if now.Sub(c.idle) >= IdleTime*time.Minute {
		return true
	}

	return false
}

// getIdle returns a string with the last time the client
// was active
func (c *Client) getIdle() string {
	if c.isIdle() {
		return fmt.Sprintf("%v", c.idle)
	}
	return "active"
}

// Connection returns the clients connection
func (c *Client) Connection() net.Conn {
	return c.conn
}

// ChannelOut returns the clients output channel
func (c *Client) ChannelOut() chan string {
	return c.output
}

// Nickname returns a string containing the Clients nickname
func (c *Client) Nickname() string {
	return c.nickname
}

// RealName returns a string containing the Clients real name
func (c *Client) Realname() string {
	return c.realname
}

// Receive provides an interface to the clients output channel
func (c *Client) Receive(str string) {
	c.output <- str
}
