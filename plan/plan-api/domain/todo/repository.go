package todo

type Repository interface {
	GetAllTodo(userName string) ([]*Todo, error)
	GetTodo(id string) (*Todo, error)
	SaveTodo(todo *Todo) error
	RemoveTodo(id string) error
}
