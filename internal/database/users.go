package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"` // when we convert the user struct to json then Password field will be removed
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (email, name, password) VALUES ($1, $2,$3)"

	// if res, err := m.DB.ExecContext(ctx, query, user.Email, user.Name, user.Password); err != nil {
	// 	fmt.Println(err)
	// 	return err
	// } else {
	// 	fmt.Println(res.RowsAffected())
	// 	fmt.Println(res.LastInsertId())
	// }
	// return nil
	row := m.DB.QueryRowContext(ctx, query, user.Email, user.Name, user.Password)
	return row.Scan(&user.ID)
}
