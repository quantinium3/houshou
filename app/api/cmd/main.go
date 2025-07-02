package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quantinium3/houshou/app/api/internal/db"
	"github.com/quantinium3/houshou/app/api/internal/handler"
)

type HError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = ErrorHandler

	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT string not found in the environment variables")
	}

	ctx := context.Background()
	conn, err := connectDatabase(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer conn.Close(ctx)
	db := db.New(conn)

	userHandler := handler.NewUserHandler(db)

	v1 := e.Group("/api/v1")
	v1.GET("/healthz", handler.GetHealth)
	v1.POST("/user", userHandler.CreateUser)

	e.Logger.Fatal(e.Start(":" + portString))
}

func connectDatabase(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message string

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprint(he.Message)

		if message == "missing or malformed jwt" {
			code = http.StatusUnauthorized
		}
	} else {
		c.Logger().Error(err)
	}

	c.JSON(code, HError{
		Status:  code,
		Message: message,
		Details: nil,
	})
}
