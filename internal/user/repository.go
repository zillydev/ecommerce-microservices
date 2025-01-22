package user

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(user *NewUserRequest) (*User, error) {
	var existingEmail string
	err := r.DB.QueryRow("SELECT email FROM users WHERE email = $1", user.Email).Scan(&existingEmail)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingEmail != "" {
		return nil, errors.New("email already exists")
	}

	query := "INSERT INTO users (name, email, preferred_notifications) VALUES ($1, $2, $3) RETURNING id"
	var userID int
	err = r.DB.QueryRow(query, user.Name, user.Email, pq.Array(user.PreferredNotifications)).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &User{ID: userID, Name: user.Name, Email: user.Email, PreferredNotifications: user.PreferredNotifications}, nil
}

func (r *Repository) UpdatePreferences(updatePreferencesRequest UpdatePreferencesRequest) error {
	query := "UPDATE users SET preferred_notifications = $1 WHERE id = $2"
	_, err := r.DB.Exec(query, pq.Array(updatePreferencesRequest.PreferredNotifications), updatePreferencesRequest.UserID)
	return err
}

func (r *Repository) GetUser(userID int) (*User, error) {
	user := &User{}
	query := "SELECT id, name, email, preferred_notifications FROM users WHERE id = $1"
	row := r.DB.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Name, &user.Email, pq.Array(&user.PreferredNotifications))
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *Repository) GetUsersByPreferences(preferences GetUsersByPreferencesRequest) ([]User, error) {
	var users []User
	query := "SELECT id, name, email, preferred_notifications FROM users WHERE preferred_notifications @> $1"
	rows, err := r.DB.Query(query, pq.Array(preferences.PreferredNotifications))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, pq.Array(&user.PreferredNotifications))
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}
