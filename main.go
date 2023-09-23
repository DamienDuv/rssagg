package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	v1Router.Get("/users", apiCfg.HandlerGetUserByAPIKey)

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