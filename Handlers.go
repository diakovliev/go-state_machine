package state_machine

import (
    "github.com/diakovliev/go-event_loop"
)

type EnterLeaveHandler struct {
    callback func(*State, event_loop.Event)
}
func NewEnterLeaveHandler(callback func(*State, event_loop.Event)) *EnterLeaveHandler {
    h := new(EnterLeaveHandler)
    h.callback = callback
    return h
}


type EventsHandler struct {
    callback func(*State, event_loop.Event) bool
}
func NewEventsHandler(callback func(*State, event_loop.Event) bool) *EventsHandler {
    h := new(EventsHandler)
    h.callback = callback
    return h
}


type IdleHandler struct {
    callback func(*State)
}
func NewIdleHandler(callback func(*State)) *IdleHandler {
    h := new(IdleHandler)
    h.callback = callback
    return h
}


type WorkHandler struct {
    callback func(*State) int
}
func NewWorkHandler(callback func(*State) int) *WorkHandler {
    h := new(WorkHandler)
    h.callback = callback
    return h
}


type StateMachineHandler struct {
    callback func(*StateMachine)
}
func NewStateMachineHandler(callback func(*StateMachine)) *StateMachineHandler {
    h := new(StateMachineHandler)
    h.callback = callback
    return h
}
