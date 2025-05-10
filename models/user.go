package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// this file provides all user related functionality

// User represents the data structure for a user.
type User struct {
	ID    uuid.UUID
	Name  string
	Email string
	// this will later include profile pictures, verifying email states, etc.

	Password []byte // passwords are hashed and salted
}

// Save saves the user into the database. It hashes and salts the password before saving.
func (u *User) Save(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `INSERT INTO CraftUser (ID, Name, Email, Password) VALUES (
		$1, $2, $3, $4
	)`, u.ID, u.Name, u.Email, u.Password)

	if err != nil {
		return err
	}

	return nil
}

// userInit initializes the CraftUser table in the database.
func userInit(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS CraftUser (
		ID uuid,
		Name text,
		Email text,
		Password bytea
	)`)

	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(ctx context.Context, conn *pgx.Conn, email string) (*User, error) {
	var user User
	err := conn.QueryRow(ctx, `SELECT ID, Name, Email, Password FROM CraftUser WHERE Email = $1`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
