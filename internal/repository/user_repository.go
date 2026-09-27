package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-api-practice/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	Update(ctx context.Context, id int, user model.User) (model.User, error)
	Delete(ctx context.Context, id int) error
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	const query = `
		SELECT id, name, email, created_at
		FROM users
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read users: %w", err)
	}

	return users, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	const query = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	const query = `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at`

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if isUniqueViolation(err) {
		return model.User{}, ErrEmailExists
	}
	if err != nil {
		return model.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, id int, user model.User) (model.User, error) {
	const query = `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
		RETURNING id, name, email, created_at`

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if isUniqueViolation(err) {
		return model.User{}, ErrEmailExists
	}
	if err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted row count: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var postgresError *pgconn.PgError
	return errors.As(err, &postgresError) && postgresError.Code == "23505"
}
