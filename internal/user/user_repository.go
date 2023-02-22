package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := `INSERT INTO users(username,email,password) VALUES(?,?,?)`
	res, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password)
	if err != nil {
		return &User{}, err
	}
	var lastInsertedId int64
	lastInsertedId, err = res.LastInsertId()
	if err != nil {
		return &User{}, err
	}
	user.ID = lastInsertedId
	return user, nil
}

func (r *repository) UserNameExist(ctx context.Context, userName string) (bool, error) {
	var count int
	query := `SELECT count(*) FROM users WHERE userName=?`
	err := r.db.QueryRowContext(ctx, query, userName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, err
}

func (r *repository) EmailExist(ctx context.Context, email string) (bool, error) {
	var count int
	query := `SELECT count(*) FROM users WHERE email=?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	query := `SELECT id,email,username,password FROM users WHERE email=?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}
