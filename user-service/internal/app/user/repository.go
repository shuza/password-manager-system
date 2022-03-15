package user

import (
	"context"
	"database/sql"
	"fmt"
	"user-service/internal/app/model"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	errUniqueViolation  = pq.ErrorCode("23505")
	insertUserQuery     = `INSERT INTO users (email, full_name, password, business_name) VALUES ($1, $2, $3, $4)`
	getUserByEmailQuery = `SELECT id, email, full_name, password, business_name FROM users WHERE email = $1`
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) InsertUser(ctx context.Context, user model.User) error {
	if _, err := r.db.ExecContext(ctx, insertUserQuery, user.Email, user.FullName, user.Password, user.BusinessName); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == errUniqueViolation {
			return fmt.Errorf("%v :%w", err, model.ErrInvalid)
		}
		return err
	}
	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	if err := r.db.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("email not found :%w", model.ErrNotFound)
		}
		return model.User{}, err
	}
	return user, nil
}
