package store

type TodoItem struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Task     string `json:"task"`
	DueDate  string `json:"dueDate"`
	Status   string `json:"status"`
	Username string `json:"username" gorm:"foreignKey:Username"`
}

type TodoStore interface {
	GetAll(username string) ([]TodoItem, error)
	GetById(username, todoId string) (TodoItem, error)
	Update(username, todoId, task, dueDate string) (TodoItem, error)
	Delete(username, id string) error
	Create(username, task, dueDate string) (TodoItem, error)
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserStore interface {
	Create(username, password, firstName, lastName string) (User, error)
	GetByUsername(username string) (User, error)
	Update(username, password, firstName, lastName string) (User, error)
	Delete(username string) error
	ValidatePassword(username, password string) (bool, error)
}

type Session struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type SessionStore interface {
	Create(username string) (Session, error)
	GetById(sessionId string) (Session, error)
	Delete(sessionId string) error
}
