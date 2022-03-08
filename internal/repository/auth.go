package repository

import (
	"JWT_auth/internal/model"
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

const result = "no rows in result set"

//save user in db
func (r *Repository) SaveUser(user *model.User) (string, error) {
	var id string
	q := `INSERT INTO users(username,phone,password,role)
    VALUES($1,$2,$3,$4)
	RETURNING id;`
	row := r.db.QueryRow(context.Background(), q, user.Username, user.Phone, user.Password, user.Role)
	row.Scan(&id)
	err := row.Scan().Error()
	if err != result {
		if strings.Contains(err, "SQLSTATE 23505") {
			return "", errors.New("error: user already exist")
		}

		return "", errors.New("error: internal DB error")
	}

	return id, nil
}

//get user from db
func (r *Repository) GetUser(user *model.User) (string, string, error) {
	var id, role string
	q := `SELECT id,role FROM users
	WHERE
		username=$1 AND password=$2;`
	r.db.QueryRow(context.Background(), q, user.Username, user.Password).Scan(&id, &role)
	if id == "" {
		return "", "", errors.New("error: user not found")
	}
	return id, role, nil
}

//get user from db
func (r *Repository) CheckUser(id uuid.UUID) (string, error) {
	var role string
	q := `SELECT role FROM users
	WHERE
		id=$1;`
	r.db.QueryRow(context.Background(), q, id).Scan(&role)
	return role, nil
}
