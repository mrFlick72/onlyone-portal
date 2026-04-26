package plan

import (
	"time"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
)

type Plan struct {
	Id       string
	UserName string
	Title    string
	Date     time.Time
}

type PlanDetails struct {
	Id    string
	Todos []*todo.Todo
}
