package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DamienDuv/rssagg/internal/database"
	"github.com/google/uuid"
)

func (server *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := server.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, user)
}

func (server *Server) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, user)
}

func (server *Server) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Limit int `json:"limit"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	posts, err := server.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(params.Limit),
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't read user's posts: %v",err))
		return
	}

	respondWithJSON(w, 200, posts)
}