package main

import (
	"fmt"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func parsePathValueUUID(r *http.Request, w http.ResponseWriter, key string) (string, uuid.UUID, bool) {
	idString := r.PathValue(key)
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid %s", key), err)
		return idString, id, false
	}
	return idString, id, true
}

func (cfg *apiConfig) validateToken(r *http.Request, w http.ResponseWriter) (uuid.UUID, bool) {

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return uuid.Nil, false
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return uuid.Nil, false
	}

	return userID, true
}

func parseMultipartForm(r *http.Request, w http.ResponseWriter, maxMemory int64) bool {
	err := r.ParseMultipartForm(maxMemory)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse multipart form", err)
		return false
	}
	return true
}

func (cfg *apiConfig) getPathIdvalidateUserIDParseMultipartForm(w http.ResponseWriter, r *http.Request, idKey string, maxMemory int64) (string, uuid.UUID, uuid.UUID, bool) {

	idString, id, ok := parsePathValueUUID(r, w, idKey)

	if !ok {
		return idString, id, uuid.Nil, false
	}

	userID, ok := cfg.validateToken(r, w)

	if !ok {
		return idString, id, userID, false
	}

	if !parseMultipartForm(r, w, maxMemory) {
		return idString, id, userID, false
	}

	return idString, id, userID, true
}
