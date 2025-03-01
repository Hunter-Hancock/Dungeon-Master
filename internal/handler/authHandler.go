package handler

import (
	"Hunter-Hancock/dungeon-master/internal/db"
	"net/http"
)

type AuthHandler struct {
	AuthStore db.AuthStore
}

func NewAuthHandler(authStore db.AuthStore) *AuthHandler {
	return &AuthHandler{AuthStore: authStore}
}

func (ah *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {

}
