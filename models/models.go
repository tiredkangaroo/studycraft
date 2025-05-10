package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// this function will be used to initialize all models.
func Init(ctx context.Context, conn *pgx.Conn) error {
	err := deckInit(ctx, conn)
	if err != nil {
		return fmt.Errorf("deck init: %w", err)
	}

	err = userInit(ctx, conn)
	if err != nil {
		return fmt.Errorf("user init: %w", err)
	}

	return nil
}
