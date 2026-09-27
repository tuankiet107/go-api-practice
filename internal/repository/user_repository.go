package repository

import (
	"sort"
	"sync"

	"go-api-practice/internal/model"
)

type UserRepository interface {
	FindAll() []model.User
	FindByID(id int) (model.User, bool)
	Create(user model.User) model.User
	Update(id int, user model.User) (model.User, bool)
	Delete(id int) bool
}

type MemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[int]model.User
	nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: map[int]model.User{
			1: {ID: 1, Name: "An", Email: "an@example.com"},
			2: {ID: 2, Name: "Binh", Email: "binh@example.com"},
		},
		nextID: 3,
	}
}

func (r *MemoryUserRepository) FindAll() []model.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users
}

func (r *MemoryUserRepository) FindByID(id int) (model.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, found := r.users[id]
	return user, found
}

func (r *MemoryUserRepository) Create(user model.User) model.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user

	return user
}

func (r *MemoryUserRepository) Update(id int, user model.User) (model.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.users[id]; !found {
		return model.User{}, false
	}

	user.ID = id
	r.users[id] = user
	return user, true
}

func (r *MemoryUserRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.users[id]; !found {
		return false
	}

	delete(r.users, id)
	return true
}
