package todo

import "time"

type Todo struct {
	Id       string
	UserName string
	Date     time.Time
	Content  string
}
