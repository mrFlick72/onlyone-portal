package plan

import (
	"time"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/logging"
)

var logger = logging.GetLoggerInstance()

type Plan struct {
	Id       string
	UserName string
	Title    string
	Date     time.Time
}

type PlanDetails struct {
	Id    string
	Todos []*Todo
}

type Todo struct {
	Id       string    `sql:"id"`
	UserName string    `sql:"user_name"`
	Date     time.Time `sql:"date"`
	Content  string    `sql:"content"`
}
