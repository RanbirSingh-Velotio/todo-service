package service

import (
	"context"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/todo"
	"github.com/RanbirSingh-Velotio/todo-service/store/sqlite"
)

type Service struct {
	store *sqlite.StoreSvc
}

func New(store *sqlite.StoreSvc) *Service {
	service := &Service{
		store: store,
	}
	return service
}

func (s *Service) TodoCreateRequest(ctx context.Context, requestInput todo.TodoRequestInput) (todo.TodoResponse, error) {
	chErr := make(chan error)
	var response todo.TodoResponse
	go func() {
		r, err := s.store.CreateTodoTask(ctx, requestInput)
		response = r
		chErr <- err
	}()

	select {
	case <-ctx.Done():
		return response, ctx.Err()
	case err := <-chErr:
		return response, err
	}
}

func (s *Service) TodoGetRequest(ctx context.Context, ids []int) []todo.TodoResponse {
	chErr := make(chan error)
	var response []todo.TodoResponse
	go func() {
		r := s.store.GetTodoTaskByID(ctx, ids)
		response = r
		chErr <- nil
	}()
	select {
	case <-ctx.Done():
		return response
	case _ = <-chErr:
		return response
	}
}
func (s *Service) TodoDeleteRequest(ctx context.Context, ids []int) todo.TodoResponse {

	response := s.store.DeleteTodoTaskByID(ctx, ids)
	return response
}
func (s *Service) TodoUpdateRequest(ctx context.Context, requestInput todo.TodoRequestInput) todo.TodoResponse {

	response := s.store.UpdateTodoTaskByID(ctx, requestInput)
	return response

}
