package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DamienDuv/rssagg/internal"
	"github.com/DamienDuv/rssagg/internal/database"
	"github.com/DamienDuv/rssagg/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatalln("DB_URL is not found in the environment")
	}

	dbConnection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	apiCfg := handlers.ApiConfig{
		DB: database.New(dbConnection),
	}

	go internal.StartScraping(apiCfg.DB, 10, time.Minute * 1)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/error", handlers.HandlerError)

	v1Router.Post("/users", apiCfg.HandlerCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))

	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Println("Server starting on port ", portString)
	errSrv := srv.ListenAndServe()
	if errSrv != nil {
		log.Fatalln(errSrv)
	}

}
