package main

import (
    "encoding/json"
	"errors"
    "net/http"
    "strings"
    "time"

	"github.com/haoq-repo/go_chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
    UserID      uuid.UUID `json:"user_id"`
    Body        string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body    string      `json:"body"`
        UserID  uuid.UUID   `json:"user_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
    }

    cleaned, err := validateChirp(params.Body)
    if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
    }

    chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
        Body:   cleaned,
        UserID: params.UserID,
    })
    respondWithJSON(w, http.StatusCreated, Chirp{
			ID:			chirp.ID,
			CreatedAt:	chirp.CreatedAt,
			UpdatedAt:	chirp.UpdatedAt,
			Body:		chirp.Body,
            UserID:     chirp.UserID,
    })
}

func validateChirp(body string) (string, error) {
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        return "", errors.New("Chirp is too long")    
    }

    cleaned := getCleanedBody(body)
    return cleaned, nil
}

func getCleanedBody(body string) string {
    profane_list := map[string]struct{}{
        "kerfuffle":    {}, 
        "sharbert":     {},
        "fornax":       {},
    }
    words := strings.Split(body, " ")

    for i, word := range words {
        loweredWord := strings.ToLower(word)
        if _, ok := profane_list[loweredWord]; ok {
            words[i] = "****"
        }
    }

    return strings.Join(words, " ")
}
