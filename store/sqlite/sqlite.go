package sqlite

import (
	"context"
	"fmt"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/todo"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type StoreSvc struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *StoreSvc {
	store := &StoreSvc{
		db: db,
	}
	return store
}

func (s *StoreSvc) GetTaskList() []todo.TodoResponse {

	queryDataSQL := "SELECT id, name, completed FROM todo"
	rows, err := s.db.Query(queryDataSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []todo.TodoResponse

	for rows.Next() {
		var todo todo.TodoResponse
		err := rows.Scan(&todo.Id, &todo.Name, &todo.Completed)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	return todos
}

func (s *StoreSvc) CreateTodoTask(ctx context.Context, requestInput todo.TodoRequestInput) (todo.TodoResponse, error) {
	_, err := s.db.Exec("PRAGMA journal_mode = WAL")
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`insert into todo(id, name,completed) values(?, ?,?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	defer tx.Rollback()
	_, err = stmt.Exec(requestInput.Id, requestInput.Name, requestInput.Completed)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return s.GetTodoTaskByID(ctx, []int{requestInput.Id})[0], nil

}

func (s *StoreSvc) GetTodoTaskByID(ctx context.Context, ids []int) []todo.TodoResponse {

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	if len(ids) == 0 {
		return s.GetTaskList()
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	// Construct the SQL query with the IN clause and placeholders.
	queryDataSQL := fmt.Sprintf("SELECT id, name, completed FROM todo WHERE id IN (%s)", strings.Join(placeholders, ","))

	// Query tasks from the 'todo' table with dynamic-length IDs.
	rows, err := tx.Query(queryDataSQL, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []todo.TodoResponse

	for rows.Next() {
		var todo todo.TodoResponse
		err := rows.Scan(&todo.Id, &todo.Name, &todo.Completed)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	return todos
}

func (s *StoreSvc) DeleteTodoTaskByID(ctx context.Context, id []int) todo.TodoResponse {
	s.db.Exec("PRAGMA journal_mode = WAL")
	deleteDataSQL := "DELETE FROM todo WHERE id = ?"
	for _, taskID := range id {
		_, err := s.db.Exec(deleteDataSQL, taskID)
		if err != nil {
			log.Printf("Error deleting task with ID %d: %v\n", taskID, err)
		} else {
			fmt.Printf("Task with ID %d deleted successfully.\n", taskID)
		}
	}
	return todo.TodoResponse{Message: "Success"}
}

func (s *StoreSvc) UpdateTodoTaskByID(ctx context.Context, requestInput todo.TodoRequestInput) todo.TodoResponse {
	s.db.Exec("PRAGMA journal_mode = WAL")
	updateDataSQL := "UPDATE todo SET name = ?, completed = ? WHERE id = ?"
	result, err := s.db.Exec(updateDataSQL, requestInput.Name, requestInput.Completed, requestInput.Id)
	if err != nil {
		log.Fatal(err)
	}

	// Check the number of rows affected by the update.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected > 0 {
		fmt.Printf("Task with ID %d updated successfully.\n", requestInput.Id)
	} else {
		fmt.Printf("No task found with ID %d.\n", requestInput.Id)
	}
	return todo.TodoResponse{
		Message:   "Success",
		Id:        requestInput.Id,
		Name:      requestInput.Name,
		Completed: requestInput.Completed,
	}
}
