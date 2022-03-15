package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"user-service/internal/app/model"

	"github.com/rs/zerolog/log"
)

func (s *server) signUp(w http.ResponseWriter, r *http.Request) {
	var data model.User
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ErrUnprocessableEntityResponse(w, "Decode Error", err)
		return
	}
	if err := s.userService.CreateUser(r.Context(), data); err != nil {
		if errors.Is(err, model.ErrInvalid) {
			ErrInvalidEntityResponse(w, "invalid user", err)
			return
		}
		log.Error().Err(err).Msgf("[signUp] failed to create user Error: %v", err)
		ErrInternalServerResponse(w, "failed to create user", err)
		return
	}
	SuccessResponse(w, http.StatusCreated, "successful")
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var data model.User
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ErrUnprocessableEntityResponse(w, "Decode Error", err)
		return
	}
	user, err := s.userService.GetUserByEmailAndPassword(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalid) || errors.Is(err, model.ErrNotFound) {
			ErrInvalidEntityResponse(w, "invalid credentials", err)
			return
		}
		log.Error().Err(err).Msgf("[login] failed to fetch login data Error: %v", err)
		ErrInternalServerResponse(w, "failed to fetch login data", err)
		return
	}

	token, err := s.authService.Encode(user)
	if err != nil {
		log.Error().Err(err).Msgf("[login] failed to generate jwt token Error: %v", err)
		ErrInternalServerResponse(w, "failed to generate jwt token", err)
		return
	}
	SuccessResponse(w, http.StatusOK, token)
}

func (s *server) tokenValidation(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	claim, err := s.authService.Decode(token)
	if err != nil {
		ErrUnauthorizedResponse(w, "invalid token", err)
		return
	}
	SuccessResponse(w, http.StatusOK, claim.User)
}
