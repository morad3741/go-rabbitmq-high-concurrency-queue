package user

import (
	"database/sql"
	"log"
)

type UserRepository interface {
	Save(user *User) error
	FindById(id string) (*User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl with a database connection.
func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save inserts a new user into the database.
func (u *UserRepositoryImpl) Save(user *User) error {
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
	_, err := u.db.Exec(query, user.ID, user.Name, user.Email)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		return err
	}
	return nil
}

// FindById retrieves a user by their ID from the database.
func (u *UserRepositoryImpl) FindById(id string) (*User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"
	row := u.db.QueryRow(query, id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found with ID: %s", id)
			return nil, nil // No user found
		}
		log.Printf("Error finding user: %v", err)
		return nil, err
	}
	return &user, nil
}
