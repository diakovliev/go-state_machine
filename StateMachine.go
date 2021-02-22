package state_machine

import (
    "fmt"
    "errors"
    "log"
    "sync"
    "time"

    "github.com/diakovliev/go-event_loop"
)

type StateMachine struct {
    Name string

    ParentSm *StateMachine
    Context interface{}
    eventLoop *event_loop.EventLoop

    onInitialize *StateMachineHandler
    onFinalize *StateMachineHandler

    states sync.Map
    finalStates map[string]bool
    transitions *TransitionsManager

    current *State

    childs []*StateMachine
}

func NewStateMachine(name string,
    parentSm *StateMachine,
    context interface{},
    onInitialize *StateMachineHandler,
    onFinalize *StateMachineHandler) *StateMachine {

    sm := new(StateMachine)

    sm.Name         = name
    sm.ParentSm     = parentSm
    sm.Context      = context
    sm.onInitialize = onInitialize
    sm.onFinalize   = onFinalize
    sm.finalStates  = make(map[string]bool)
    sm.transitions  = NewTransitionsManager()
    sm.current      = nil
    sm.childs       = make([]*StateMachine, 0, 100)

    sm.setupEventLoop()

    if sm.ParentSm != nil {
        sm.childs = append(sm.childs, sm)
    }

    return sm
}

func (sm *StateMachine) setupEventLoop() {

    sm.eventLoop = event_loop.NewEventLoop(sm.Name + "@events")

    sm.eventLoop.SetDefaultHandler(event_loop.NewDefaultEventsHandler(sm.Name + "%default", func(event event_loop.Event) bool {
        sm.handleEvent(event, nil)
        return true
    }))

    sm.eventLoop.SetLoopHandler(event_loop.NewNoEventsHandler(sm.Name + "%loop", func() {
        if (sm.current.relaxCPU) {
            time.Sleep(100 * time.Millisecond)
            sm.current.relaxCPU = false
        }

        if sm.current.WorkComplete {
            sm.current.callOnIdle()
            sm.current.relaxCPU = true

        } else  {
            ret := sm.current.callOnWork()

            sm.current.WorkComplete     = (ret & STATE_WORK_COMPLETE) != 0
            sm.current.relaxCPU         = (ret & STATE_WORK_RELAX_CPU) != 0
        }
    }))
}

func (sm *StateMachine) handleEvent(event event_loop.Event, innerMachine *StateMachine) {
    t := sm.findTransition(event)
    if t != nil {

        if (innerMachine != nil) {

            log.Printf("%s: Inner machine '%s': Leave state '%s'", sm.Name, innerMachine.Name, sm.current.Name)
            innerMachine.current.callOnLeave(event)

            log.Printf("%s: Stop inner machine '%s'", sm.Name, innerMachine.Name)
            innerMachine.Stop()

            log.Printf("%s: Resent event from event machine.", sm.Name)
            sm.Send(event_loop.NewEvent(event.Type, event.Payload))
            return
        }

        log.Printf("%s: Call goto state %s.", sm.Name, t.to)
        sm.performTransition(t.to, event)

    } else {

        log.Printf("%s: No transition. Call state event handler.", sm.Name)

        handlingComplete := sm.current.callOnEvents(event)
        if handlingComplete {
            log.Printf("%s: Event handling complete.", sm.Name)
            return
        }

        if sm.ParentSm != nil {
            log.Printf("%s: Not completed event handling. Handle event by parent machine.", sm.Name)
            sm.ParentSm.handleEvent(event, sm)
        }
    }
}

func (sm *StateMachine) findTransition(event event_loop.Event) *Transition {

    log.Printf("%s: Find transition by event '%d'", sm.Name, event.Type)

    t := sm.transitions.Find(sm.current.Name, event)
    if t != nil {
        log.Printf("%s: Found transition %s -> %s", sm.Name, t.from, t.to)
        return t
    }

    log.Printf("%s: Find global transition by event '%d'", sm.Name, event.Type)
    t = sm.transitions.FindGlobal(event)
    if t != nil {
        log.Printf("%s: Found global transition %s -> %s", sm.Name, t.from, t.to)
        return t
    }

    log.Printf("%s: No transitions by event '%d'", sm.Name, event.Type)
    return nil
}

func (sm *StateMachine) getState(state_name string) *State {
    istate, ok := sm.states.Load(state_name)
    if !ok {
        panic(fmt.Errorf("Can't find state '%s'!", state_name))
    }

    state, ok := istate.(*State)
    if !ok {
        panic(errors.New("Not a *State in states!"))
    }

    return state
}

func (sm *StateMachine) performTransition(state_name string, event event_loop.Event) {
    log.Printf("%s: Leave state '%s'", sm.Name, sm.current.Name)
    sm.current.callOnLeave(event)

    sm.current = sm.getState(state_name)

    log.Printf("%s: Enter state '%s'", sm.Name, sm.current.Name)
    sm.current.callOnEnter(event)

    _, final := sm.finalStates[sm.current.Name]
    if final {
        log.Printf("%s: State '%s' is final. Stop machine event loop.", sm.Name, sm.current.Name)
        sm.eventLoop.Stop(0)
    }
}

func (sm *StateMachine) addState(state *State) {
    _, ok := sm.states.Load(state.Name)
    if ok {
        panic(fmt.Errorf("%s: Non unique state name '%s'!", sm.Name, state.Name))
    }

    log.Printf("%s: Add state '%s'.", sm.Name, state.Name)
    sm.states.Store(state.Name, state)
}

func (sm *StateMachine) callOnInitialize() {
    if sm.onInitialize != nil {
        sm.onInitialize.callback(sm)
    }
}

func (sm *StateMachine) callOnFinalize() {
    if sm.onFinalize != nil {
        sm.onFinalize.callback(sm)
    }
}

func (sm *StateMachine) CreateInitialState(name string,
    onEvents *EventsHandler,
    onLeave *EnterLeaveHandler,
    onWork *WorkHandler,
    onIdle *IdleHandler) *State {

    state := newState(sm, name, nil, onEvents, onLeave, onWork, onIdle)

    sm.addState(state)

    log.Printf("%s: Set initial state '%s'.", sm.Name, state.Name)
    sm.current = state

    return state
}
func (sm *StateMachine) CreateState(name string,
    onEnter *EnterLeaveHandler,
    onEvents *EventsHandler,
    onLeave *EnterLeaveHandler,
    onWork *WorkHandler,
    onIdle *IdleHandler) *State {

    state := newState(sm, name, onEnter, onEvents, onLeave, onWork, onIdle)

    sm.addState(state)

    return state
}
func (sm *StateMachine) CreateFinalState(name string,
    onEnter *EnterLeaveHandler) *State {

    state := newState(sm, name, onEnter, nil, nil, nil, nil)

    sm.addState(state)

    log.Printf("%s: Add state '%s' to final states set.", sm.Name, state.Name)
    sm.finalStates[name] = true

    return state
}


func (sm *StateMachine) AddTransition(from *State, to *State, condition OnEventCondition) {
    log.Printf("%s: Add %s -> %s transition.", sm.Name, from.Name, to.Name)
    sm.transitions.Add(from.Name, to.Name, condition)
}

func (sm *StateMachine) AddGlobalTransition(to *State, condition OnEventCondition) {
    log.Printf("%s: Add * -> %s transition.", sm.Name, to.Name)
    sm.transitions.AddGlobal(to.Name, condition)
}

func (sm *StateMachine) Run() {
    log.Printf("%s: Initialize.", sm.Name)
    sm.callOnInitialize()

    defer func() {
        log.Printf("%s: Finalize.", sm.Name)
        sm.callOnFinalize()
    }()

    log.Printf("%s: Run event loop.", sm.Name)
    sm.eventLoop.Run()
}

func (sm *StateMachine) Send(event *event_loop.Event) {
    log.Printf("%s: Put event '%d' to event loop.", sm.Name, event.Type)
    sm.eventLoop.Send(event)
}

func (sm *StateMachine) Stop() {
    for _, child := range(sm.childs) {
        child.Stop()
    }
    log.Printf("%s: Stop machine event loop.", sm.Name)
    sm.eventLoop.Stop(0)
}
