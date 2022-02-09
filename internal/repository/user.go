package repository

import (
	"context"
	"fmt"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Login          string `json:"login" db:"login"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
}

func (r *Repository) Login(ctx context.Context, login, hashedPassword string) (u User, err error) {
	row := r.pool.QueryRow(ctx, `select id, login, name, surname from users where login = $1 AND hashed_password = $2`, login, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}
