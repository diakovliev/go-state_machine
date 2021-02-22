package state_machine

import (
    "github.com/diakovliev/go-event_loop"
)

type TransitionsManager struct {
    transitions []*Transition
    globalTransitions []*Transition
}

func NewTransitionsManager() *TransitionsManager {
    tm := new(TransitionsManager)
    tm.transitions          = make([]*Transition, 0, 100)
    tm.globalTransitions    = make([]*Transition, 0, 100)
    return tm
}

func (tm *TransitionsManager) Add(from string, to string, condition OnEventCondition) {
    transition := new(Transition)
    transition.from         = from
    transition.to           = to
    transition.condition    = condition

    tm.transitions = append(tm.transitions, transition)
}

func (tm *TransitionsManager) AddGlobal(to string, condition OnEventCondition) {
    transition := new(Transition)
    transition.from         = "*"
    transition.to           = to
    transition.condition    = condition

    tm.globalTransitions = append(tm.globalTransitions, transition)
}

func (tm *TransitionsManager) Find(state string, event event_loop.Event) *Transition {
    for _, t := range(tm.transitions) {
        if t.from == state && t.condition.condition(event) {
            return t
        }
    }
    return nil
}

func (tm *TransitionsManager) FindGlobal(event event_loop.Event) *Transition {
    for _, t := range(tm.globalTransitions) {
        if t.condition.condition(event) {
            return t
        }
    }
    return nil
}
