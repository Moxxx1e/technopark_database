package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/technopark_database/internal/delivery"
	"github.com/technopark_database/internal/user/repository"
	"github.com/technopark_database/internal/user/usecases"
	"log"
)

func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "technodb", "123456", "technodb")
}

func main() {
	e := echo.New()

	db, err := sql.Open("postgres", GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserPgRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo)
	userHandler := delivery.NewUserHandler(userUseCase)
	userHandler.Configure(e)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
