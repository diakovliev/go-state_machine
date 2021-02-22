package state_machine

import (
    "github.com/diakovliev/go-event_loop"
)

type OnEventCondition interface {
    condition(event_loop.Event) bool
}



type OnAnyEventType struct {}
func (c* OnAnyEventType) condition(event_loop.Event) bool {
    return true
}
func OnAnyEvent() *OnAnyEventType {
    return new(OnAnyEventType)
}



type OnEventTypeType struct {
    Type uint64
}
func (c* OnEventTypeType) condition(event event_loop.Event) bool {
    return event.Type == c.Type
}
func OnEventType(expectedType uint64) *OnEventTypeType {
    c := new(OnEventTypeType)
    c.Type = expectedType
    return c
}



type Transition struct {
    from string
    to string

    condition OnEventCondition
}
