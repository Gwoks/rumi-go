package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"rumi-go/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewUserStore(db *sqlx.DB, timeout time.Duration) *UserStore {
	return &UserStore{
		db:      db,
		timeout: timeout,
	}
}

var (
	getUserQuery        = `SELECT id, email, password, name, phone, role, address FROM users WHERE id = ?`
	getUserByEmailQuery = `SELECT id, email, password, name, phone, role, address FROM users WHERE email = ?`
	createUserQuery     = `INSERT INTO users (email, password, name, phone, role, address) VALUES (?, ?, ?, ?, ?, ?)`
	updateUserQuery     = `UPDATE users SET email = ?, password = ?, name = ?, phone = ?, role = ?, address = ? WHERE id = ?`
	deleteUserQuery     = `DELETE FROM users WHERE id = ?`
	listUsersQuery      = `SELECT id, email, password, name, phone, role, address FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	countUsersQuery     = `SELECT COUNT(*) FROM users`
)

func (s *UserStore) Get(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User

	err := s.db.GetContext(ctx, &user, getUserQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := s.db.GetContext(ctx, &user, getUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, createUserQuery,
		user.Email, user.Password, user.Name, user.Phone, user.Role, user.Address)
	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (s *UserStore) Update(ctx context.Context, user *model.User) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.db.ExecContext(ctxWithDeadline, updateUserQuery,
		user.Email, user.Password, user.Name, user.Phone, user.Role, user.Address, user.ID)
	return err
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, deleteUserQuery, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStore) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var users []*model.User
	err := s.db.SelectContext(ctxWithDeadline, &users, listUsersQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStore) Count(ctx context.Context) (int64, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var count int64
	err := s.db.GetContext(ctxWithDeadline, &count, countUsersQuery)
	return count, err
}
