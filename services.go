package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RouterTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Primeira Api com GO rodando",
	})
}

// BUSCANDO TODAS AS TAREFAS
func GetAllTasks(c *gin.Context) {
	rows, err := DB.Query("SELECT id, title FROM tasks")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	defer rows.Close()

	var tasks []Tasks

	for rows.Next() {
		var task Tasks
		if err := rows.Scan(&task.Id, &task.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}

// CADASTRAR NOVA TAREFA
func AddNewTask(c *gin.Context) {
	var newTask Tasks

	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := DB.Exec("INSERT INTO tasks (title) VALUES (?)", newTask.Title)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newTask.Id = int(id)
	c.JSON(http.StatusCreated, newTask)

}

// BUSCANDO TAREFA PELO ID
func GetTaskById(c *gin.Context) {
	id := c.Param("id")

	var task Tasks
	row := DB.QueryRow("SELECT id, title FROM tasks WHERE id = ?", id)

	if err := row.Scan(&task.Id, &task.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DELETANDO PELO ID
func DeleteTaskById(c *gin.Context) {
	id := c.Param("id")

	_, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarefa deletada com sucesso!"})
}

// ATUALIZAR PELO ID
func UpdateTaskById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updateTask Tasks

	if err := c.BindJSON(&updateTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tarefa não encontrada!"})
		return
	}

	_, err := DB.Exec("UPDATE tasks SET title = ? WHERE id = ?", updateTask.Title, id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	updateTask.Id = id
	c.JSON(http.StatusOK, updateTask)
}
