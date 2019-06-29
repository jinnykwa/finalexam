package todo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jinnykwa/finalexam/database"
)

type Todo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Todohandler struct{}

func (Todohandler) PostTodosHandler(c *gin.Context) {
	fmt.Println("PostTodosHandler")
	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(t)

	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	query := `INSERT INTO custom (name, email, status) VALUES ($1,$2,$3) RETURNING id`
	var id int
	row := db.QueryRow(query, t.Name, t.Email, t.Status)
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	t.ID = id
	fmt.Println("Insert success id :", id)
	c.JSON(201, t)
}

func (Todohandler) GetTodosHandler(c *gin.Context) {
	fmt.Println("GetTodosHandler")
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, email, status FROM custom WHERE id = $1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	myid := c.Param("id")

	row := stmt.QueryRow(myid)
	t := Todo{}

	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Select one row is", t.ID, t.Name, t.Email, t.Status)
	c.JSON(http.StatusOK, t)
}

func (Todohandler) GetlistTodosHandler(c *gin.Context) {
	fmt.Println("GetlistTodosHandler")
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, email, status FROM custom")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	custom := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		custom = append(custom, t)
	}
	c.JSON(http.StatusOK, custom)
}

func (Todohandler) PutupdateTodosHandler(c *gin.Context) {
	fmt.Println("PutupdateTodosHandler")
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE custom SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	myid := c.Param("id")

	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	t.ID, err = strconv.Atoi(myid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if _, err := stmt.Exec(myid, t.Name, t.Email, t.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exec error": err.Error()})
		return
	}
	fmt.Println("Update success", t.ID, t.Name, t.Email, t.Status)
	c.JSON(http.StatusOK, t)
}

func (Todohandler) DeleteTodosByIdHandler(c *gin.Context) {
	fmt.Println("DeleteTodosByIdHandler")
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM custom WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	myid := c.Param("id")

	stmt.QueryRow(myid)
	fmt.Println("Delete success")
	c.JSON(http.StatusOK, gin.H{
		"message": "customer deleted",
	})
}
