package internal

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/unmsmfisi-socialapplication/social_app/pkg/database"

	loginApplication "github.com/unmsmfisi-socialapplication/social_app/internal/login/application"
	loginInfrastructure "github.com/unmsmfisi-socialapplication/social_app/internal/login/infrastructure"
	"github.com/unmsmfisi-socialapplication/social_app/internal/post_reactions/handler"
	"github.com/unmsmfisi-socialapplication/social_app/internal/post_reactions/repository"
	"github.com/unmsmfisi-socialapplication/social_app/internal/post_reactions/service"
	// postReactionApplication "github.com/unmsmfisi-socialapplication/social_app/internal/post_reactions/application"
	// postReactionInfrastructure "github.com/unmsmfisi-socialapplication/social_app/internal/post_reactions/infrastructure"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	dbInstance := database.GetDB()

	dbRepo := loginInfrastructure.NewUserDBRepository(dbInstance)
	loginUseCase := loginApplication.NewLoginUseCase(dbRepo)
	loginHandler := loginInfrastructure.NewLoginHandler(loginUseCase)

	postReactionRepository := repository.NewPostReactionRepository(dbInstance)
	postReactionService := service.NewPostReactionService(postReactionRepository)
	postReactionHandler := handler.NewPostReactionHandler(postReactionService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		// Simulates some hard work.
		//
		// We want this handler to complete successfully during a shutdown signal,
		// so consider the work here as some background routine to fetch a long running
		// search query to find as many results as possible, but, instead we cut it short
		// and respond with what we have so far. How a shutdown is handled is entirely
		// up to the developer, as some code blocks are preemptible, and others are not.
		time.Sleep(5 * time.Second)

		w.Write([]byte(fmt.Sprintf("{\"response\": \"all done slow\"}")))
	})

	// Login
	r.Post("/login", loginHandler.HandleLogin)

	r.Get("/reactions/{post_id}", postReactionHandler.GetReactionsForPost)
	r.Post("/reactions", postReactionHandler.CreatePostReactionHandler)

	return r
}
