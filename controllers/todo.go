package controllers

import (
	"example/go-htmx/store"
	"example/go-htmx/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todoRouter struct {
	todoStore store.TodoStore
	renderer  views.Renderer
}

type TodoRouterParams struct {
	TodoStore store.TodoStore
	Renderer  views.Renderer
}

func NewTodoRouter(params TodoRouterParams) *todoRouter {
	return &todoRouter{
		todoStore: params.TodoStore,
		renderer:  params.Renderer,
	}
}

func (r *todoRouter) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", r.getTodoPage)
	router.GET("/item/:id", r.getTodoItem)
	router.GET("/edititem/:id", r.getEditTodoItem)
	router.POST("/item", r.postTodoItem)
	router.PUT("/item/:id", r.putTodoItem)
	router.DELETE("/item/:id", r.deleteTodoItem)
}

func (r *todoRouter) getTodoPage(c *gin.Context) {
	user, _ := store.GetGlobals(c).GetUser()

	items, err := r.todoStore.GetAll(user.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.NewRenderer(c).AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}
	views.NewRenderer(c).TodoPage(items).Render(c, c.Writer)
}

func (r *todoRouter) getTodoItem(c *gin.Context) {
	user, exists := store.GetGlobals(c).GetUser()

	if !exists {
		c.Status(http.StatusUnauthorized)
		c.Header("HX-Redirect", "/auth/login")
		return
	}

	id := c.Param("id")

	if id == "" {
		c.Status(http.StatusBadRequest)
		views.NewRenderer(c).AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	item, err := r.todoStore.GetById(user.Username, id)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.NewRenderer(c).AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.NewRenderer(c).TodoItem(item).Render(c, c.Writer)
}

func (r *todoRouter) getEditTodoItem(c *gin.Context) {

	user := store.GetGlobals(c).User

	// Get the ID from the URL parameter
	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		views.NewRenderer(c).AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	item, error := r.todoStore.GetById(user.Username, id)

	if error != nil {
		c.Status(http.StatusInternalServerError)
		views.NewRenderer(c).AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.NewRenderer(c).EditTodoItem(item).Render(c, c.Writer)
}

func (r *todoRouter) postTodoItem(c *gin.Context) {

	c.Request.ParseForm()
	formData := c.Request.PostForm
	task := formData.Get("task")
	dueDate := formData.Get("due_date")

	if task == "" || dueDate == "" {
		c.Status(http.StatusBadRequest)
		views.NewRenderer(c).AddItemError("Task and Due Date are required fields.").Render(c, c.Writer)
		return
	}

	user := store.GetGlobals(c).User

	newItem, err := r.todoStore.Create(user.Username, task, dueDate)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.NewRenderer(c).AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.NewRenderer(c).TodoItem(newItem).Render(c, c.Writer)
}

func (r *todoRouter) putTodoItem(c *gin.Context) {

	todoId := c.Param("id")

	if todoId == "" {
		c.Status(http.StatusBadRequest)
		views.NewRenderer(c).AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	c.Request.ParseForm()
	formData := c.Request.PostForm
	task := formData.Get("task")
	dueDate := formData.Get("due_date")

	if task == "" || dueDate == "" {
		c.Status(http.StatusBadRequest)
		views.NewRenderer(c).AddItemError("Task and Due Date are required fields.").Render(c, c.Writer)
		return
	}

	user := store.GetGlobals(c).User

	newItem, err := r.todoStore.Update(user.Username, todoId, task, dueDate)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		views.NewRenderer(c).AddItemError("Error. Please try again later.").Render(c, c.Writer)
		return
	}

	views.NewRenderer(c).TodoItem(newItem).Render(c, c.Writer)
}

func (r *todoRouter) deleteTodoItem(c *gin.Context) {

	// Get the ID from the URL parameter
	id := c.Param("id")

	if id == "" {
		c.Status(http.StatusBadRequest)
		views.NewRenderer(c).AddItemError("ID is required.").Render(c, c.Writer)
		return
	}

	user := store.GetGlobals(c).User

	r.todoStore.Delete(user.Username, id)
}
