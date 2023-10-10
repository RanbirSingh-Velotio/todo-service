package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/httputil"
	"github.com/RanbirSingh-Velotio/todo-service/pkg/todo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	service todo.Service
}

var (
	errBadRequest     = errors.New("BAD_REQUEST")
	errRequestTimeOut = errors.New("REQUEST_TIMEOUT")
)

func InitHandler(service todo.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetIdentity returns handler identity
func (h *Handler) GetIdentity() string {
	return "todo-v1"
}

// Start will start all http handlers
func (h *Handler) Start() error {
	http.Handle("/v1/todo", TraceMiddleware(h))
	return nil
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) errorResponse(w http.ResponseWriter, code int) {

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.HandleCreateRequest(w, r)
	case http.MethodGet:
		h.HandleGetRequest(w, r)
	case http.MethodPut:
		h.HandlePutRequest(w, r)
	case http.MethodDelete:
		h.HandleDeleteRequest(w, r)
	default:
		// Return error immediately if the request method is incorrect
		h.errorResponse(w, http.StatusMethodNotAllowed)
	}
}

func (h *Handler) parseTodoQueryParam(ctx context.Context, r *http.Request) ([]int, error) {

	idsParam := r.URL.Query().Get("ids")

	if idsParam == "" {
		return []int{}, nil
	}

	// Split the comma-separated string into individual IDs.
	idStrings := strings.Split(idsParam, ",")

	var ids []int
	for _, idStr := range idStrings {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err // Handle parsing errors
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (h *Handler) parseTodoRequest(ctx context.Context, r *http.Request) (todo.TodoRequestInput, error) {
	var inputRequest todo.TodoRequestInput
	var err error
	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		// Parse request
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errChan <- errBadRequest
			return
		}

		err = json.Unmarshal(body, &inputRequest)
		if err != nil {
			errChan <- errBadRequest
			return
		}
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		err = errRequestTimeOut
		return inputRequest, nil
	case err = <-errChan:
		if err != nil {
			return inputRequest, err
		}
	}

	return inputRequest, nil
}

func (h *Handler) HandleCreateRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var inputRequestData todo.TodoRequestInput
	var err error
	errChan := make(chan error, 1)
	var response todo.TodoResponse
	defer func(start time.Time) {
		jsonResponse, _ := json.Marshal(response)
		_, err := httputil.WriteResponse(w, jsonResponse, http.StatusOK, httputil.NewContentTypeDecorator("application/json"))
		if err != nil {
			return
		}
	}(time.Now())

	go func(ctx context.Context) {
		inputRequestData, err = h.parseTodoRequest(ctx, r)
		if err != nil {
			errChan <- err
			return
		}

		response, err = h.service.TodoCreateRequest(ctx, inputRequestData)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		return
	}

}

func (h *Handler) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	errChan := make(chan error, 1)
	var response []todo.TodoResponse
	defer func(start time.Time) {
		jsonResponse, _ := json.Marshal(response)
		_, err := httputil.WriteResponse(w, jsonResponse, http.StatusOK, httputil.NewContentTypeDecorator("application/json"))
		if err != nil {
			return
		}
	}(time.Now())
	go func(ctx context.Context) {
		var ids []int
		ids, err = h.parseTodoQueryParam(ctx, r)
		if err != nil {
			errChan <- err
			return
		}

		response = h.service.TodoGetRequest(ctx, ids)
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		return
	}
}

func (h *Handler) HandlePutRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var inputRequestData todo.TodoRequestInput
	var err error
	errChan := make(chan error, 1)
	var response todo.TodoResponse
	defer func(start time.Time) {
		jsonResponse, _ := json.Marshal(response)
		_, err := httputil.WriteResponse(w, jsonResponse, http.StatusOK, httputil.NewContentTypeDecorator("application/json"))
		if err != nil {
			return
		}
	}(time.Now())

	go func(ctx context.Context) {
		inputRequestData, err = h.parseTodoRequest(ctx, r)
		if err != nil {
			errChan <- err
			return
		}

		response = h.service.TodoUpdateRequest(ctx, inputRequestData)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		return
	}
}

func (h *Handler) HandleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	errChan := make(chan error, 1)
	var response todo.TodoResponse
	defer func(start time.Time) {
		jsonResponse, _ := json.Marshal(response)
		_, err := httputil.WriteResponse(w, jsonResponse, http.StatusOK, httputil.NewContentTypeDecorator("application/json"))
		if err != nil {
			return
		}
	}(time.Now())
	go func(ctx context.Context) {
		var ids []int
		ids, err = h.parseTodoQueryParam(ctx, r)
		if err != nil {
			errChan <- err
			return
		}

		response = h.service.TodoDeleteRequest(ctx, ids)
		errChan <- nil
	}(ctx)

	select {
	case <-ctx.Done():
		return
	case err = <-errChan:
		return
	}
}
