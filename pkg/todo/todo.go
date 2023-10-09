package todo

import (
	"context"
)

type TodoRequestInput struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type TodoResponse struct {
	Message   string `json:"message"`
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Completed bool   `json:"completed,omitempty"`
}

//go:generate mockgen -destination mockservice/mock_service.go -package mockservice github.com/RanbirSingh-Velotio/todo-service/pkg/todo Service
type Service interface {
	TodoCreateRequest(ctx context.Context, requestInput TodoRequestInput) (TodoResponse, error)
	TodoGetRequest(ctx context.Context, ids []int) []TodoResponse
	TodoDeleteRequest(ctx context.Context, ids []int) TodoResponse
	TodoUpdateRequest(ctx context.Context, input TodoRequestInput) TodoResponse
}

var defaultService Service

func Init(svc Service) {
	defaultService = svc
}

func GetService() Service {
	return defaultService
}
