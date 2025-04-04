package controllers

import (
	"example/go-htmx/middleware"
	"example/go-htmx/request"
	"example/go-htmx/store"
	"example/go-htmx/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todoRouter struct {
	todoStore store.TodoStore
	userStore store.UserStore
}

type TodoRouterParams struct {
	TodoStore store.TodoStore
	UserStore store.UserStore
}

func NewTodoRouter(params TodoRouterParams) *todoRouter {
	if params.TodoStore == nil {
		panic("TodoStore is required")
	} else if params.UserStore == nil {
		panic("UserStore is required")
	}
	return &todoRouter{
		todoStore: params.TodoStore,
		userStore: params.UserStore,
	}
}

func (r *todoRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.getTodoPage)
	router.GET("/item/:id", middleware.RequireAuth(r.getTodoItem))
	router.GET("/edititem/:id", middleware.RequireAuth(r.getEditTodoItem))
	router.POST("/item", middleware.RequireAuth(r.postTodoItem))
	router.PUT("/item/:id", middleware.RequireAuth(r.putTodoItem))
	router.DELETE("/item/:id", middleware.RequireAuth(r.deleteTodoItem))
}

func (r *todoRouter) getTodoPage(c *gin.Context) {
	user, exists := request.GetUser(c, r.userStore)
	var items []store.TodoItem

	if exists {
		items, _ = r.todoStore.GetAll(user.Username)
	}

	views.TodoPage(user, items).Render(c, c.Writer)
}

func (r *todoRouter) getTodoItem(c *gin.Context) {
	username, _ := request.GetUsername(c)

	id := c.Param("id")

	if id == "" {
		c.Status(http.StatusBadRequest)
		views.AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	item, err := r.todoStore.GetById(username, id)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.TodoItem(item).Render(c, c.Writer)
}

func (r *todoRouter) getEditTodoItem(c *gin.Context) {
	username, _ := request.GetUsername(c)

	// Get the ID from the URL parameter
	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		views.AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	item, error := r.todoStore.GetById(username, id)

	if error != nil {
		c.Status(http.StatusInternalServerError)
		views.AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.EditTodoItem(item).Render(c, c.Writer)
}

func (r *todoRouter) postTodoItem(c *gin.Context) {

	c.Request.ParseForm()
	formData := c.Request.PostForm
	task := formData.Get("task")
	dueDate := formData.Get("due_date")

	if task == "" || dueDate == "" {
		c.Status(http.StatusBadRequest)
		views.AddItemError("Task and Due Date are required fields.").Render(c, c.Writer)
		return
	}

	username, _ := request.GetUsername(c)

	newItem, err := r.todoStore.Create(username, task, dueDate)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.TodoItem(newItem).Render(c, c.Writer)
}

func (r *todoRouter) putTodoItem(c *gin.Context) {

	todoId := c.Param("id")

	if todoId == "" {
		c.Status(http.StatusBadRequest)
		views.AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	c.Request.ParseForm()
	formData := c.Request.PostForm
	task := formData.Get("task")
	dueDate := formData.Get("due_date")

	if task == "" || dueDate == "" {
		c.Status(http.StatusBadRequest)
		views.AddItemError("Task and Due Date are required fields.").Render(c, c.Writer)
		return
	}

	username, _ := request.GetUsername(c)

	newItem, err := r.todoStore.Update(username, todoId, task, dueDate)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.TodoItem(newItem).Render(c, c.Writer)
}

func (r *todoRouter) deleteTodoItem(c *gin.Context) {

	// Get the ID from the URL parameter
	id := c.Param("id")

	if id == "" {
		c.Status(http.StatusBadRequest)
		views.AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	username, _ := request.GetUsername(c)

	r.todoStore.Delete(username, id)
}
