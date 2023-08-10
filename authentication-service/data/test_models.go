package data

import (
	"database/sql"
	"time"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{Conn: db}
}

func (r *PostgresTestRepository) GetAll() ([]*User, error) {
	var users []*User

	return users, nil
}

func (r *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "admin@admin.com",
		FirstName: "First",
		LastName:  "Last",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (r *PostgresTestRepository) GetOne(id int) (*User, error) {
	user := User{
		ID:        1,
		Email:     "admin@admin.com",
		FirstName: "First",
		LastName:  "Last",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (r *PostgresTestRepository) Update(user User) error {
	return nil
}

func (r *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

func (r *PostgresTestRepository) Insert(user User) (int, error) {
	return 2, nil
}

func (r *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

func (r *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}
