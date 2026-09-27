package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"go-api-practice/internal/model"
	"go-api-practice/internal/repository"
)

func TestHealth(t *testing.T) {
	response := sendRequest(newRouter(newFakeUserRepository()), http.MethodGet, "/health", nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}
	if response.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestUserCRUD(t *testing.T) {
	r := newRouter(newFakeUserRepository())

	t.Run("list initial users", func(t *testing.T) {
		response := sendRequest(r, http.MethodGet, "/users", nil)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
		}

		var users []model.User
		if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
			t.Fatal(err)
		}
		if len(users) != 2 {
			t.Fatalf("expected 2 initial users, got %d", len(users))
		}
	})

	t.Run("get user by id", func(t *testing.T) {
		response := sendRequest(r, http.MethodGet, "/users/1", nil)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
		}
	})

	t.Run("reject invalid user", func(t *testing.T) {
		body := []byte(`{"name":"","email":"not-an-email"}`)
		response := sendRequest(r, http.MethodPost, "/users", body)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
		}
	})

	t.Run("create update and delete user", func(t *testing.T) {
		createBody := []byte(`{"name":"Chi","email":"chi@example.com"}`)
		response := sendRequest(r, http.MethodPost, "/users", createBody)
		if response.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
		}

		var created model.User
		if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
			t.Fatal(err)
		}

		updateBody := []byte(`{"name":"Chi Updated","email":"chi.updated@example.com"}`)
		response = sendRequest(r, http.MethodPut, "/users/3", updateBody)
		if response.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
		}

		response = sendRequest(r, http.MethodDelete, "/users/3", nil)
		if response.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
		}

		response = sendRequest(r, http.MethodGet, "/users/3", nil)
		if response.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
		}
	})
}

func sendRequest(handler http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

type fakeUserRepository struct {
	users  map[int]model.User
	nextID int
}

func newFakeUserRepository() *fakeUserRepository {
	createdAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	return &fakeUserRepository{
		users: map[int]model.User{
			1: {ID: 1, Name: "Mia", Email: "mia@example.com", CreatedAt: createdAt},
			2: {ID: 2, Name: "Lisa", Email: "lisa@example.com", CreatedAt: createdAt},
		},
		nextID: 3,
	}
}

func (r *fakeUserRepository) FindAll(_ context.Context) ([]model.User, error) {
	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	return users, nil
}

func (r *fakeUserRepository) FindByID(_ context.Context, id int) (model.User, error) {
	user, found := r.users[id]
	if !found {
		return model.User{}, repository.ErrUserNotFound
	}
	return user, nil
}

func (r *fakeUserRepository) Create(_ context.Context, user model.User) (model.User, error) {
	user.ID = r.nextID
	user.CreatedAt = time.Now()
	r.nextID++
	r.users[user.ID] = user
	return user, nil
}

func (r *fakeUserRepository) Update(_ context.Context, id int, user model.User) (model.User, error) {
	existing, found := r.users[id]
	if !found {
		return model.User{}, repository.ErrUserNotFound
	}
	user.ID = id
	user.CreatedAt = existing.CreatedAt
	r.users[id] = user
	return user, nil
}

func (r *fakeUserRepository) Delete(_ context.Context, id int) error {
	if _, found := r.users[id]; !found {
		return repository.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}
