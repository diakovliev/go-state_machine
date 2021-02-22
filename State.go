package state_machine

import (
    "github.com/diakovliev/go-event_loop"
)

const (
    STATE_WORK_COMPLETE     = 0x01
    STATE_WORK_RELAX_CPU    = 0x02
    STATE_WORK_CONTINUE     = 0x00
)

type State struct {
    Sm              *StateMachine

    Name            string
    WorkComplete    bool
    relaxCPU        bool

    onEnter         *EnterLeaveHandler
    onEvents        *EventsHandler
    onLeave         *EnterLeaveHandler
    onIdle          *IdleHandler
    onWork          *WorkHandler
}

func newState(sm *StateMachine,
    name string,
    onEnter *EnterLeaveHandler,
    onEvents *EventsHandler,
    onLeave *EnterLeaveHandler,
    onWork *WorkHandler,
    onIdle *IdleHandler) *State {

    s := new(State)

    s.Sm            = sm
    s.Name          = name

    s.onEnter       = onEnter
    s.onEvents      = onEvents
    s.onLeave       = onLeave
    s.onWork        = onWork
    s.onIdle        = onIdle

    s.WorkComplete  = false
    s.relaxCPU      = false

    return s
}

func (s *State) callOnWork() int {
    if s.onWork == nil {
        return STATE_WORK_COMPLETE
    } else {
        return s.onWork.callback(s)
    }
}

func (s *State) callOnIdle() {
    if s.onIdle != nil {
        s.onIdle.callback(s)
    }
}

func (s *State) callOnEvents(event event_loop.Event) bool {
    if s.onEvents == nil {
        return false
    } else {
        return s.onEvents.callback(s, event)
    }
}

func (s *State) callOnEnter(event event_loop.Event) {
    if s.onEnter != nil {
        s.onEnter.callback(s, event)
    }
}

func (s *State) callOnLeave(event event_loop.Event) {
    if s.onLeave != nil {
        s.onLeave.callback(s, event)
    }
    s.WorkComplete  = false
    s.relaxCPU      = false
}
