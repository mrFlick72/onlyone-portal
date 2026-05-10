package plan

import (
	"time"
)

type Plan struct {
	Id       string
	UserName string
	Title    string
	Date     time.Time
	Todos    []*Todo
}

type Todo struct {
	Id       string
	UserName string
	Date     time.Time
	Content  string
	Status   TodoStatus
}

type TodoStatus string

const (
	StatusTodo       TodoStatus = "TODO"
	StatusInProgress TodoStatus = "IN_PROGRESS"
	StatusDone       TodoStatus = "DONE"
	StatusAborted    TodoStatus = "ABORTED"
)

var allowedTransitions = map[TodoStatus]map[TodoStatus]struct{}{
	StatusTodo:       {StatusInProgress: {}, StatusAborted: {}},
	StatusInProgress: {StatusDone: {}, StatusAborted: {}},
	StatusDone:       {},
	StatusAborted:    {},
}

func (s TodoStatus) IsValid() bool {
	_, ok := allowedTransitions[s]
	return ok
}

func (s TodoStatus) CanTransitionTo(target TodoStatus) bool {
	targets, ok := allowedTransitions[s]
	if !ok {
		return false
	}
	_, ok = targets[target]
	return ok
}
