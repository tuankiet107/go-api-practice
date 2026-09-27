package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-practice/internal/model"
	"go-api-practice/internal/router"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.New().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	expectedBody := "{\"status\":\"ok\"}\n"
	if recorder.Body.String() != expectedBody {
		t.Fatalf("expected body %q, got %q", expectedBody, recorder.Body.String())
	}
}

func TestUserCRUD(t *testing.T) {
	r := router.New()

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

		var user model.User
		if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
			t.Fatal(err)
		}
		if user.ID != 1 || user.Name != "An" {
			t.Fatalf("unexpected user: %+v", user)
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
		if created.ID != 3 {
			t.Fatalf("expected new user ID 3, got %d", created.ID)
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
