package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id                  string
	Email               string
	Name                string
	AvatarUrl           string
	RefreshTokenVersion int
}

type Identity struct {
	UserId        string
	Provider      string
	ProviderEmail string
	ProviderId    string
}

type AuthStore interface {
	CreateUser(user *User) error
	CreateIdentity(identity *Identity) error
	Exists(providerEmail, providerId string) *User
	ProviderExists(identity *Identity) bool
	GetUser(id string) (*User, error)
	IncrementTokenVersion(userId string) error
}

type SqlAuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *SqlAuthStore {
	return &SqlAuthStore{db: db}
}

func (s SqlAuthStore) CreateUser(user *User) error {
	query := "INSERT INTO user (id, email, avatar_url, refresh_token_version) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, user.Id, user.Email, user.AvatarUrl, user.RefreshTokenVersion)
	if err != nil {
		fmt.Println("Error inserting user", err)
		return err
	}

	return nil
}

func (s SqlAuthStore) CreateIdentity(identity *Identity) error {
	query := "INSERT INTO identity (user_id, provider, provider_email, provider_id) VALUES (?, ?, ?, ?)"

	_, err := s.db.Exec(query, identity.UserId, identity.Provider, identity.ProviderEmail, identity.ProviderId)
	if err != nil {
		fmt.Println("Error inserting identity", err)
		return err
	}

	return nil
}

func (s SqlAuthStore) Exists(providerEmail, providerId string) *User {
	var user User
	row := s.db.QueryRow("SELECT id, email, avatar_url, refresh_token_version from user WHERE email = ?", providerEmail)
	err := row.Scan(&user.Id, &user.Email, &user.AvatarUrl, &user.RefreshTokenVersion)
	if err == nil {
		return &user
	}

	row = s.db.QueryRow("SELECT user_id FROM identity WHERE provider_email = ? AND provider_id = ?", providerEmail, providerId)
	err = row.Scan(&user.Id)
	if err != nil {
		fmt.Println("Error checking user", err)
		return nil
	}
	return &user
}

func (s SqlAuthStore) ProviderExists(identity *Identity) bool {
	var provider string
	row := s.db.QueryRow("SELECT provider from identity where user_id = ? AND provider = ?", identity.UserId, identity.Provider)
	err := row.Scan(&provider)
	if err != nil {
		return false
	}

	return true
}

func (s SqlAuthStore) GetUser(id string) (*User, error) {
	query := "SELECT id, email, avatar_url FROM user WHERE id = ?"
	row := s.db.QueryRow(query, id)
	user := &User{}
	err := row.Scan(&user.Id, &user.Email, &user.AvatarUrl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s SqlAuthStore) IncrementTokenVersion(userId string) error {
	var currVersion int
	row := s.db.QueryRow("SELECT refresh_token_version FROM user where id = ?", userId)
	err := row.Scan(&currVersion)
	if err != nil {
		return err
	}
	currVersion++
	_, err = s.db.Exec("UPDATE user SET refresh_token_version = ?", currVersion)
	if err != nil {
		return err
	}

	return nil
}
