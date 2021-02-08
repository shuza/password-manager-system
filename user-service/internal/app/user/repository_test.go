package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"user-service/internal/app/model"
)

func TestRepository_InsertUser(t *testing.T) {
	t.Run("should return success", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		defer db.Close()

		user := model.User{
			Email:        "email",
			FullName:     "full_name",
			Password:     "password",
			BusinessName: "business_name",
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		m.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs(user.Email, user.FullName, user.Password, user.BusinessName).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewRepository(sqlxDB)
		err := repo.InsertUser(context.Background(), user)
		assert.Nil(t, err)
	})

	t.Run("should return unique key violation error", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		defer db.Close()

		user := model.User{
			Email:        "email",
			FullName:     "full_name",
			Password:     "password",
			BusinessName: "business_name",
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		m.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs(user.Email, user.FullName, user.Password, user.BusinessName).
			WillReturnError(&pq.Error{Code: "23505"})

		repo := NewRepository(sqlxDB)
		err := repo.InsertUser(context.Background(), user)
		assert.True(t, errors.Is(err, model.ErrInvalid))
	})

	t.Run("should return sql error", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		defer db.Close()

		user := model.User{
			Email:        "email",
			FullName:     "full_name",
			Password:     "password",
			BusinessName: "business_name",
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		m.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs(user.Email, user.FullName, user.Password, user.BusinessName).
			WillReturnError(errors.New("sql-error"))

		repo := NewRepository(sqlxDB)
		err := repo.InsertUser(context.Background(), user)
		assert.NotNil(t, err)
	})
}

func TestRepository_GetUserByEmail(t *testing.T) {
	t.Run("should return success", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		defer db.Close()

		user := model.User{
			ID:           1,
			Email:        "abc@gmail.com",
			FullName:     "mr abc",
			Password:     "123456",
			BusinessName: "business-1",
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		m.ExpectQuery("^SELECT (.+) FROM users WHERE (.+)").
			WithArgs("email").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "full_name", "password", "business_name"}).
				AddRow(1, "abc@gmail.com", "mr abc", "123456", "business-1"))
		repo := NewRepository(sqlxDB)
		result, err := repo.GetUserByEmail(context.Background(), "email")
		assert.Nil(t, err)
		assert.EqualValues(t, user, result)
	})

	t.Run("should return no rows error", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		m.ExpectQuery("^SELECT (.+) FROM users WHERE (.+)").
			WithArgs("email").
			WillReturnError(sql.ErrNoRows)
		repo := NewRepository(sqlxDB)
		_, err := repo.GetUserByEmail(context.Background(), "email")
		assert.True(t, errors.Is(err, model.ErrNotFound))
	})

	t.Run("should return error", func(t *testing.T) {
		db, m, _ := sqlmock.New()
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		m.ExpectQuery("^SELECT (.+) FROM users WHERE (.+)").
			WithArgs("email").
			WillReturnError(errors.New("sql-error"))
		repo := NewRepository(sqlxDB)
		_, err := repo.GetUserByEmail(context.Background(), "email")
		assert.NotNil(t, err)
	})
}
