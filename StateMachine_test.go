package state_machine

import (
    "log"
    "testing"

    "github.com/diakovliev/go-event_loop"
)

const (
    EVENT_ERROR = event_loop.EVENT_PRE_USER + 1
    EVENT1      = event_loop.EVENT_PRE_USER + 2
)

func TestCreateStateMachine(t *testing.T) {

    onWork := NewWorkHandler(func(state *State) int {
        if state.Name == "s2" {
            log.Print("Send EVENT_ERROR")
            state.Sm.Send(event_loop.NewEvent(EVENT_ERROR, nil))
            return STATE_WORK_COMPLETE
        } else {
            log.Print("Send EVENT1")
            state.Sm.Send(event_loop.NewEvent(EVENT1, nil))
            return STATE_WORK_COMPLETE
        }
    })

    sm := NewStateMachine("test_sm", nil, nil, nil, nil)
    if sm == nil {
        t.Fatal("Can't create state machine!")
    }

    s1 := sm.CreateInitialState("s1",    nil, nil, onWork, nil)
    if s1 == nil {
        t.Fatal("Can't create state!")
    }
    s2 := sm.CreateState("s2",           nil, nil, nil, onWork, nil)
    if s2 == nil {
        t.Fatal("Can't create state!")
    }
    s3 := sm.CreateState("s3",           nil, nil, nil, onWork, nil)
    if s3 == nil {
        t.Fatal("Can't create state!")
    }
    s4 := sm.CreateFinalState("s4",      nil)
    if s4 == nil {
        t.Fatal("Can't create state!")
    }
    sE := sm.CreateFinalState("error",   nil)
    if sE == nil {
        t.Fatal("Can't create state!")
    }

    sm.AddTransition(s1, s2, OnEventType(EVENT1))
    sm.AddTransition(s2, s3, OnEventType(EVENT1))
    sm.AddTransition(s3, s4, OnEventType(EVENT1))

    sm.AddGlobalTransition(sE, OnEventType(EVENT_ERROR))

    sm.Run()
}
