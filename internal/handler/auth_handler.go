package handler

import (
	"Hunter-Hancock/dungeon-master/internal/db"
	"Hunter-Hancock/dungeon-master/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"net/http"
)

type AuthHandler struct {
	db db.AuthStore
}

func NewAuthHandler(authStore db.AuthStore) *AuthHandler {
	return &AuthHandler{db: authStore}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, req)
	if err != nil {
		fmt.Println("Authentication Failed: ", err)
		return
	}

	existingUser := ah.db.Exists(user.Email, user.UserID)
	var userId string

	if existingUser == nil {
		newUser := &db.User{
			Id:                  uuid.New().String(),
			Email:               user.Email,
			AvatarUrl:           user.AvatarURL,
			RefreshTokenVersion: 1,
		}

		fmt.Println(newUser)

		err = ah.db.CreateUser(newUser)
		if err != nil {
			fmt.Println("Error creating user", err)
			return
		}

		identity := &db.Identity{
			UserId:        newUser.Id,
			Provider:      user.Provider,
			ProviderId:    user.UserID,
			ProviderEmail: user.Email,
		}

		err = ah.db.CreateIdentity(identity)
		if err != nil {
			fmt.Println("Error creating identity", err)
		}
		userId = newUser.Id
	} else {
		userId = existingUser.Id
		identity := &db.Identity{
			UserId:        userId,
			Provider:      user.Provider,
			ProviderId:    user.UserID,
			ProviderEmail: user.Email,
		}
		if !ah.db.ProviderExists(identity) {
			err := ah.db.CreateIdentity(identity)
			if err != nil {
				fmt.Println("Error adding identity", err)
				return
			}
		}
	}

	auth.SendTokenCookies(w, existingUser)
	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearCookies(w)

	w.Header().Set("Location", "http://localhost:5173/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (ah *AuthHandler) Begin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))

	if _, err := gothic.CompleteUserAuth(w, req); err == nil {
	} else {
		gothic.BeginAuthHandler(w, req)
	}
}

func (ah *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromReq(r)
	if err != nil {
		fmt.Println("[HANDLER] Error getting user id", err)
		return
	}

	user, err := ah.db.GetUser(userId)
	if err != nil {
		fmt.Println("[HANDLER] Error getting user", err)
		return
	}

	auth.SendTokenCookies(w, user)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		fmt.Println("[HANDLER] Error encoding user", err)
	}
}
