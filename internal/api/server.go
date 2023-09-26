package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/DamienDuv/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Server struct {
	DB database.Querier
}

func NewServer(dbUrl string) *Server {
	dbConnection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{
		DB: database.New(dbConnection),
	}
}

func (s *Server) StartListening(portString string) error {
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
	v1Router.Get("/healthz", HandlerReadiness)
	v1Router.Get("/error", HandlerError)

	v1Router.Post("/users", s.HandlerCreateUser)
	v1Router.Get("/users", s.MiddlewareAuth(s.HandlerGetUser))
	v1Router.Get("/users/posts", s.MiddlewareAuth(s.handlerGetPostsForUser))

	v1Router.Post("/feeds", s.MiddlewareAuth(s.HandlerCreateFeed))
	v1Router.Get("/feeds", s.HandlerGetFeeds)

	v1Router.Post("/feed_follows", s.MiddlewareAuth(s.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", s.MiddlewareAuth(s.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", s.MiddlewareAuth(s.HandlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)
	
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	return srv.ListenAndServe()
}

func (s *Server) StartScraping(concurrency int, interval time.Duration)  {
	go StartScraping(s.DB, concurrency, interval)
}
	