package wst

import "fmt"

type Event struct {
	Target IEventDispatcher
	Type   string
	Object interface{}
}

type EventDispatcher struct {
	savers []*EventSaver
}

type EventSaver struct {
	Type      string
	Listeners []*EventListener
}

type EventListener struct {
	Handler EventHandler
}

type EventHandler func(event Event)

type IEventDispatcher interface {
	AddEventListener(eventType string, listener *EventListener)
	RemoveEventListener(eventType string, listener *EventListener) bool
	HasEventListener(eventType string) bool
	DispatchEvent(event Event) bool
}

func NewEventDispatcher() *EventDispatcher {
	return new(EventDispatcher)
}

func NewEventListener(h EventHandler) *EventListener {
	l := new(EventListener)
	l.Handler = h
	return l
}

func NewEvent(eventType string, object interface{}) Event {
	e := Event{Type: eventType, Object: object}
	return e
}

func (this *Event) Clone() *Event {
	e := new(Event)
	e.Type = this.Type
	e.Target = e.Target
	return e
}

func (this *Event) ToString() string {
	return fmt.Sprintf("Event Type %v", this.Type)
}

func (this *EventDispatcher) AddEventListener(eventType string, listener *EventListener) {
	for _, saver := range this.savers {
		if saver.Type == eventType {
			saver.Listeners = append(saver.Listeners, listener)
			return
		}
	}

	saver := &EventSaver{Type: eventType, Listeners: []*EventListener{listener}}
	this.savers = append(this.savers, saver)
}

func (this *EventDispatcher) RemoveEventListener(eventType string, listener *EventListener) bool {
	for _, saver := range this.savers {
		if saver.Type == eventType {
			for i, l := range saver.Listeners {
				if listener == l {
					saver.Listeners = append(saver.Listeners[:i], saver.Listeners[i+1:]...)
					return true
				}
			}
		}
	}
	return false
}

func (this *EventDispatcher) HasEventListener(eventType string) bool {
	for _, saver := range this.savers {
		if saver.Type == eventType {
			return true
		}
	}
	return false
}

func (this *EventDispatcher) DispatchEvent(event Event) bool {
	for _, saver := range this.savers {
		if saver.Type == event.Type {
			for _, listener := range saver.Listeners {
				event.Target = this
				listener.Handler(event)
			}
			return true
		}
	}
	return false
}
