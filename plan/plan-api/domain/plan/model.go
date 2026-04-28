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
}
