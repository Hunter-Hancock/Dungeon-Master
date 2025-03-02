package auth

import (
	"Hunter-Hancock/dungeon-master/internal/db"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

func CreateAuthTokens(user *db.User) *TokenResponse {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":              user.Id,
		"refreshTokenVersion": user.RefreshTokenVersion,
		"exp":                 time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})

	refreshTokenString, _ := refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	accessTokenString, _ := accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}
}

func ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	id, ok := token.Claims.(jwt.MapClaims)["userId"]
	if ok {
		return id.(string), nil
	}

	return "", fmt.Errorf("token is invalid")
}

func ParseRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	id, ok := token.Claims.(jwt.MapClaims)["userId"]
	if ok {
		return id.(string), nil
	}

	return "", fmt.Errorf("token is invalid")
}

func GetUserIdFromReq(r *http.Request) (string, error) {
	accessTokenCookie, err := r.Cookie("id")
	if err != nil {
		fmt.Println("No Access Token, we need to check refresh and send new one")
	} else {
		accessTokenUserId, err := ParseAccessToken(accessTokenCookie.Value)
		if err != nil {
			fmt.Println("Error getting user id", err)
		}

		if accessTokenUserId != "" {
			return accessTokenUserId, nil
		}
	}

	refreshTokenCookie, err := r.Cookie("rid")
	if err != nil {
		return "", err
	}

	refreshTokenUserId, err := ParseRefreshToken(refreshTokenCookie.Value)
	if err != nil {
		fmt.Println("Error getting user id", err)
		return "", err
	}

	if refreshTokenUserId != "" {
		return refreshTokenUserId, nil
	}

	return "", fmt.Errorf("token is invalid")
}

func SendTokenCookies(res http.ResponseWriter, user *db.User) {
	tokens := CreateAuthTokens(user)
	accessCookie := &http.Cookie{
		Name:     "id",
		Value:    tokens.AccessToken,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(time.Minute * 15),
	}

	refreshCookie := &http.Cookie{
		Name:     "rid",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
	}

	http.SetCookie(res, accessCookie)
	http.SetCookie(res, refreshCookie)
}

func ClearCookies(res http.ResponseWriter) {
	accessCookie := &http.Cookie{
		Name:     "id",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now(),
	}

	refreshCookie := &http.Cookie{
		Name:     "rid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now(),
	}
	http.SetCookie(res, accessCookie)
	http.SetCookie(res, refreshCookie)
}
