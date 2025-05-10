package main

import (
	"main/models"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// choosing echo for now - not because of familiarity but instead because of reliability
// and ease of use

func Serve(conn *pgx.Conn) error {
	e := echo.New()
	// use echo logging

	e.POST("/api/signup", func(c echo.Context) error {
		err := c.Request().ParseMultipartForm(1e4) // 10kb limit should help avoid resource exhaustion
		if err != nil {
			return jsonerror(c, 500, "failed to read request body")
		}

		id, err := uuid.NewRandom()
		if err != nil {
			return jsonerror(c, 500, "failed to generate user ID")
		}

		username := c.Request().FormValue("username")
		password := c.Request().FormValue("password")
		email := c.Request().FormValue("email")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return jsonerror(c, 500, "failed to hash password, user not created")
		}

		user := models.User{
			ID:       id,
			Name:     username,
			Password: []byte(hashedPassword),
			Email:    email,
		}

		err = user.Save(c.Request().Context(), conn)
		if err != nil {
			return jsonerror(c, 500, "failed to save user")
		}

		return c.JSON(200, map[string]string{
			"id": id.String(),
		})
	})

	e.Server.Addr = os.Getenv("API_ADDR")
	return e.Server.ListenAndServe()
}

func jsonerror(c echo.Context, status int, error string) error {
	return c.JSON(status, map[string]string{
		"error": error,
	})
}

// func jsonsuccess(c echo.Context, status int) error {
// 	return c.JSON(status, map[string]string{
// 		"success": "true",
// 	})
// }
