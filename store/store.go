package store

import (
	"context"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/todo"
)

//go:generate mockgen -destination mockservice/mock_service.go -package mockservice github.com/RanbirSingh-Velotio/todo-service/pkg/store Service
type StoreSvc interface {
	GetTaskList() []todo.TodoResponse
	CreateTodoTask(ctx context.Context, requestInput todo.TodoRequestInput) (todo.TodoResponse, error)
	GetTodoTaskByID(ctx context.Context, id []int) []todo.TodoResponse
	DeleteTodoTaskByID(ctx context.Context, id []int) todo.TodoResponse
	UpdateTodoTaskByID(ctx context.Context, requestInput todo.TodoRequestInput) todo.TodoResponse
}

var defaultService StoreSvc

func Init(svc StoreSvc) {
	defaultService = svc
}

func GetService() StoreSvc {
	return defaultService
}
