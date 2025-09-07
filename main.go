package main

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

/*
A simplistic API that allows for CRUD operations on a list of TODOs.
You can:
- Get all TODOs
- Get a TODO by ID
- Add a new TODO
- Delete a TODO by ID
- Mark a TODO as done or not done

The data is stored in memory for simplicity.
*/
type Todo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type RequestContext = gin.Context

var todos = []Todo{
	{ID: "1", Name: "Learn Go", Done: false},
	{ID: "2", Name: "Build a web app", Done: false},
	{ID: "3", Name: "Destroy a web app", Done: true},
	{ID: "4", Name: "Build an API in GO", Done: false},
}

func getTodos(c *RequestContext) {
	c.IndentedJSON(http.StatusOK, todos)
}

func getTodoById(id string) (*Todo, error) {
	for i, todo := range todos {
		if todo.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("Todo not found")
}

func getTodoByIdHandler(c *RequestContext) {
	id := c.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, todo)
}

func addTodo(c *RequestContext) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		return
	}
	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func deleteTodo(id string) error {
	for i, todo := range todos {

		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			return nil
		}
	}
	return errors.New("Todo doesn't exist")
}

func deleteTodoHandler(c *RequestContext) {
	id := c.Param("id")

	err := deleteTodo(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(http.StatusNoContent)
	}

}

func markAsDone(c *RequestContext) {
	id := c.Param("id")
	var updates map[string]any
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			if done, ok := updates["done"].(bool); ok {
				todos[i].Done = done
			}
			c.JSON(http.StatusOK, todos[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodoByIdHandler)
	router.POST("/todos", addTodo)
	router.DELETE("/todos/:id", deleteTodoHandler)
	router.PATCH("todos/:id", markAsDone)
	router.Run("localhost:8080")
}
