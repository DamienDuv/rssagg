package handlers

import (
	"fmt"
	"net/http"

	"github.com/DamienDuv/rssagg/internal/auth"
	"github.com/DamienDuv/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
