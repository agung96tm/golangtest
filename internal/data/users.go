package data

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	RoleID     string    `json:"role_id"`
	RoleName   string    `json:"role_name"`
	LastAccess time.Time `json:"last_access"`
}

func (u *User) Matches(password string) bool {
	return u.Password == password
}

func (u *UserModel) GetAll() ([]User, error) {
	stmt := `SELECT id, name, email, password, role_id, role_name, last_access FROM users`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleName, &user.LastAccess); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserModel) GetByID(id int) (*User, error) {
	stmt := `SELECT id, name, email, password, role_id, role_name, last_access FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := u.DB.QueryRowContext(ctx, stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.RoleName)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	stmt := `SELECT id, name, password, role_id, role_name, last_access FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := u.DB.QueryRowContext(ctx, stmt, email).Scan(&user.ID, &user.Name, &user.Password, &user.RoleID, &user.RoleName, &user.LastAccess)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *UserModel) Create(user *User) error {
	stmt := `
		INSERT INTO users (id, name, email, password, role_id, role_name, last_access) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, stmt, user.ID, user.Name, user.Email, user.Password, user.RoleID, user.RoleName).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Update(user *User) error {
	stmt := `UPDATE 
    	SET name = $1, email = $2, password = $3, last_access = $4, role_id = $5, role_name = $6
		WHERE id = $7
		RETURNING id
	`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password,
		user.LastAccess,
		user.RoleID,
		user.RoleName,
		user.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.ID)
}

func (u *UserModel) DeleteByID(id int) error {
	stmt := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
