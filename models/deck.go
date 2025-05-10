package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Deck represents the data structure for a deck of cards.
type Deck struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string

	CreatedOn    time.Time
	LastModified time.Time
}

// Save saves the deck into the database.
func (d *Deck) Save(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `INSERT INTO Deck (ID, Name, Description, CreatedOn, LastModified) VALUES (
		$1, $2, $3, $4, $5
	)`, d.ID, d.Name, d.Description, d.CreatedOn, d.LastModified) // specifies columns for my sanity 😔🥀

	if err != nil {
		return err
	}

	return nil
}

// deckInit initialized the Deck table in the database.
func deckInit(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS Deck (
		ID uuid,
		Name text,
		Description text,
		CreatedOn timestamp,
		LastModified timestamp
	)`)
	if err != nil {
		return err
	}
	return nil
}

func GetDecksByUserID(ctx context.Context, conn *pgx.Conn, userID uuid.UUID) ([]Deck, error) {
	rows, err := conn.Query(ctx, `SELECT ID, Name, Description, CreatedOn, LastModified FROM Deck WHERE UserID = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []Deck
	for rows.Next() {
		var d Deck
		if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.CreatedOn, &d.LastModified); err != nil {
			return nil, err
		}
		decks = append(decks, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return decks, nil
}
