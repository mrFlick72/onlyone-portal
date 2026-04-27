package plan

import (
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"time"
)

type Plan struct {
	Id       string
	UserName string
	Title    string
	Date     time.Time
	Todos    []*todo.Todo
}
