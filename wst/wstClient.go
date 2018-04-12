package wst

import (
	"errors"
	"io"
	"log"
	"time"
)

const maxQueueMsgCount = 1024

type wstClient struct {
	id string
	// rwc is the interface to access the websocket connection.
	// It is set after the client registers with the server.
	rwc io.ReadWriteCloser
	// msgs is the queued message sent from this client.
	msgs []string
	// timer is used to remove this client if unregistered after a timeout.
	timer *time.Timer
}

func newClient(id string, t *time.Timer) *wstClient {
	c := wstClient{id: id, timer: t}
	return &c
}

func (c *wstClient) setTimer(t *time.Timer) {
	if c.timer != nil {
		c.timer.Stop()
	}
	c.timer = t
}

// register binds the ReadWriteCloser to the client if it's not done yet.
func (c *wstClient) register(rwc io.ReadWriteCloser) error {
	if c.rwc != nil {
		log.Printf("Not registering because the client %s already has a connection.", c.id)
		return errors.New("Duplicated registration")
	}
	c.setTimer(nil)
	c.rwc = rwc
	return nil
}

// deregister closes the ReadWriteCloser if it exists.
func (c *wstClient) deregister() {
	if c.rwc != nil {
		c.rwc.Close()
		c.rwc = nil
	}
}

// registered returns true if the client has registered.
func (c *wstClient) registered() bool {
	return c.rwc != nil
}

// enqueue adds a message to the client's message queue.
func (c *wstClient) enqueue(msg string) error {
	if len(c.msgs) >= maxQueueMsgCount {
		return errors.New("Too many messages queued for the client")
	}
	c.msgs = append(c.msgs, msg)
	return nil
}

// sendQueued the queued message to the other client.
func (c *wstClient) sendQueued(other *wstClient) error {
	if c.id == other.id || other.rwc == nil {
		return errors.New("Invalid client")
	}
	for _, m := range c.msgs {
		sendServerMsg(other.rwc, m)
	}
	c.msgs = nil
	log.Printf("Sent queued message from %s to %s", c.id, other.id)
	return nil
}

// send sends the message to the other client if the other client has registered,
// or queues the message otherwise.
func (c *wstClient) send(other *wstClient, msg string) error {
	if c.id == other.id {
		return errors.New("Invalied client")
	}
	if other.rwc != nil {
		return sendServerMsg(other.rwc, msg)
	}
	return c.enqueue(msg)
}
