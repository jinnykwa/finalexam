package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinnykwa/finalexam/todo"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Cannot connect database")
		return
	}
	defer db.Close()

	createTb := `
    CREATE TABLE IF NOT EXISTS custom (
	    id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		panic("Cannot create database")
		return
	}
	fmt.Println("Okay")
	r := gin.Default()

	s := todo.Todohandler{}
	r.POST("customers", s.PostTodosHandler)
	r.GET("customers/:id", s.GetTodosHandler)
	r.GET("customers", s.GetlistTodosHandler)
	r.PUT("customers/:id", s.PutupdateTodosHandler)
	r.DELETE("customers/:id", s.DeleteTodosByIdHandler)
	r.Run(":2019")
}
