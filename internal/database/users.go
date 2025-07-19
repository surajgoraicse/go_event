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

// func (m *UserModel) Get(id int) (*User, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := "SELECT * FORM users where id = $1"
// 	user := &User{}
// 	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
// 	if err != nil {
// 		return nil, err

// 	}
// 	return user, nil

// }

// func (m *UserModel) GetUserByEmail(email string) (*User, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := "SELECT * FROM users where email = $1"
// 	user := &User{}
// 	if err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("getting user by email : ", user)
// 	return user, nil
// }

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err

	}
	return user, nil

}

func (m *UserModel) Get(id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	return m.getUser(query, id)
}
func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	return m.getUser(query, email)
}
