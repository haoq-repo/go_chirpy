package main

import (
	"net/http"

	"github.com/haoq-repo/go_chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
        return
    }

    userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
        return
    }

	dbChirp, err := cfg.db.GetChirp(r.Context(),chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return		
	}


	err = cfg.db.DeleteChirp(r.Context(),chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chrip", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}