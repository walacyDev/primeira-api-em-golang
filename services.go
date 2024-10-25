package main

import (
	"fmt"
	"net/http"

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

	for _, task := range taskList {
		if fmt.Sprintf("%d", task.Id) == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Tarefa não foi encontrada com esse id",
	})
}

// DELETANDO PELO ID
func DeleteTaskById(c *gin.Context) {
	id := c.Param("id")

	for index, task := range taskList {
		if fmt.Sprintf("%d", task.Id) == id {
			taskList = append(taskList[:index], taskList[index+1:]...)
			c.JSON(http.StatusOK, gin.H{"Message": "Tarefa deletada com sucesso!"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Tarefa não foi encontrada",
	})
}

// ATUALIZAR PELO ID
func UpdateTaskById(c *gin.Context) {
	id := c.Param("id")

	var updateTask Tasks

	if err := c.BindJSON(&updateTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tarefa não encontrada!"})
		return
	}

	for index, task := range taskList {
		if fmt.Sprintf("%d", task.Id) == id {
			updateTask.Id = task.Id
			taskList[index] = updateTask
			c.JSON(http.StatusOK, gin.H{"message": "Tarefa atualizada com sucesso!"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Tarefa não encontrada!"})
}
