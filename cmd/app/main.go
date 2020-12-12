package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	userDelivery "github.com/technopark_database/internal/user/delivery"
	userRepository "github.com/technopark_database/internal/user/repository"
	userUseCase "github.com/technopark_database/internal/user/usecases"

	voteRepository "github.com/technopark_database/internal/vote/repository"
	voteUseCase "github.com/technopark_database/internal/vote/usecases"

	forumDelivery "github.com/technopark_database/internal/forum/delivery"
	forumRepository "github.com/technopark_database/internal/forum/repository"
	forumUseCase "github.com/technopark_database/internal/forum/usecases"

	postDelivery "github.com/technopark_database/internal/post/delivery"
	postRepository "github.com/technopark_database/internal/post/repository"
	postUseCase "github.com/technopark_database/internal/post/usecases"

	serviceDelivery "github.com/technopark_database/internal/service/delivery"
	serviceRepository "github.com/technopark_database/internal/service/repository"
	serviceUseCase "github.com/technopark_database/internal/service/usecases"

	threadDelivery "github.com/technopark_database/internal/thread/delivery"
	threadRepository "github.com/technopark_database/internal/thread/repository"
	threadUseCase "github.com/technopark_database/internal/thread/usecases"

	"log"
)

func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "techno_db")
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

	// User
	userRepo := userRepository.NewUserPgRepository(db)
	userUseCase := userUseCase.NewUserUseCase(userRepo)
	userHandler := userDelivery.NewUserHandler(userUseCase)

	// Forum
	forumRepo := forumRepository.NewForumPgRepository(db)
	forumUseCase := forumUseCase.NewForumUseCase(forumRepo, userUseCase)
	forumHandler := forumDelivery.NewForumHandler(forumUseCase)

	// Vote
	voteRepo := voteRepository.NewVoteRepository(db)
	voteUseCase := voteUseCase.NewVoteUseCase(voteRepo)

	// Thread
	threadRepo := threadRepository.NewThreadPgRepository(db)
	threadUseCase := threadUseCase.NewThreadUseCase(threadRepo, userUseCase, forumUseCase, voteUseCase)
	threadHandler := threadDelivery.NewThreadHandler(threadUseCase)

	// Service
	serviceRepo := serviceRepository.NewServicePgRepository(db)
	serviceUseCase := serviceUseCase.NewServiceUseCase(serviceRepo)
	serviceHandler := serviceDelivery.NewServiceHandler(serviceUseCase)

	postRepo := postRepository.NewPostPgRepository(db)
	postUseCase := postUseCase.NewPostUseCase(postRepo)
	postHandler := postDelivery.NewPostHandler(postUseCase)

	userHandler.Configure(e)
	forumHandler.Configure(e)
	serviceHandler.Configure(e)
	threadHandler.Configure(e)
	postHandler.Configure(e)

	e.Logger.Fatal(e.Start("localhost:5000"))
}
