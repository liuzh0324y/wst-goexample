package wst

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const maxRoomCapacity = 2

type wstRoom struct {
	parent          *wstRoomTable
	id              string
	clients         map[string]*wstClient
	registerTimeOut time.Duration
	roomSrcUrl      string
}

func newRoom(p *wstRoomTable, id string, to time.Duration, rs string) *wstRoom {
	return &wstRoom{parent: p, id: id, clients: make(map[string]*wstClient), registerTimeOut: to, roomSrcUrl: rs}
}

func (rm *wstRoom) client(clientID string) (*wstClient, error) {
	if c, ok := rm.clients[clientID]; ok {
		return c, nil
	}
	if len(rm.clients) >= maxRoomCapacity {
		log.Printf("Room %s is full, not adding client %s", rm.id, clientID)
		return nil, errors.New("Max room capacity reached")
	}

	var timer *time.Timer
	if rm.parent != nil {
		timer = time.AfterFunc(rm.registerTimeOut, func() {
			if c := rm.clients[clientID]; c != nil {
				rm.parent.removeIfUnregistered(rm.id, c)
			}
		})
	}
	rm.clients[clientID] = newClient(clientID, timer)

	log.Printf("Added client %s to room %s", clientID, rm.id)

	return rm.clients[clientID], nil
}

func (rm *wstRoom) register(clientID string, rwc io.ReadWriteCloser) error {
	c, err := rm.client(clientID)
	if err != nil {
		return err
	}
	if err = c.register(rwc); err != nil {
		return err
	}

	log.Printf("Client %s registered in room %s", clientID, rm.id)

	if len(rm.clients) > 1 {
		for _, otherClient := range rm.clients {
			otherClient.sendQueued(c)
		}
	}
	return nil
}

func (rm *wstRoom) send(srcClientID string, msg string) error {
	src, err := rm.client(srcClientID)
	if err != nil {
		return err
	}

	if len(rm.clients) == 1 {
		return rm.clients[srcClientID].enqueue(msg)
	}

	for _, oc := range rm.clients {
		if oc.id != srcClientID {
			return src.send(oc, msg)
		}
	}

	return errors.New(fmt.Sprintf("Corrupted room %+v", rm))
}

func (rm *wstRoom) remove(clientID string) {
	if c, ok := rm.clients[clientID]; ok {
		c.deregister()
		delete(rm.clients, clientID)
		log.Printf("Removed client %s from room %s", clientID, rm.id)

		resp, err := http.Post(rm.roomSrcUrl+"/bye/"+rm.id+"/"+clientID, "text", nil)
		if err != nil {
			log.Printf("Failed to post BYE to room server %s: %v", rm.roomSrcUrl, err)
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}

func (rm *wstRoom) empty() bool {
	return len(rm.clients) == 0
}

func (rm *wstRoom) wsCount() int {
	count := 0
	for _, c := range rm.clients {
		if c.registered() {
			count += 1
		}
	}
	return count
}
