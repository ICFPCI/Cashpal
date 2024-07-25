package middleware

import (
	"cashpal/api/utils"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserContextKey = ContextKey("UserID")

type Middleware func(http.Handler) http.Handler

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(middlewares); i++ {
			next = middlewares[i](next)
		}
		return next
	}
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWritter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.Split(r.Header.Get("Authorization"), " ")

		if len(tokenString) < 2 {
			http.Error(w, "authorization header not provided", http.StatusUnauthorized)
			return
		}

		if tokenString[0] != "Bearer" {
			http.Error(w, "bearer token not provided.", http.StatusUnauthorized)
			return
		}

		token, err := utils.VerifyAccessToken(tokenString[1])

		if err != nil {
			http.Error(w, "the provided jwt token is invalid or has expired. please check the token and try again.", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			log.Println("jwt.Claims cannot be parsed to jwt.MapClaims")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		userID, ok := claims["user_id"].(float64)

		if !ok {
			log.Println("user_id cannot be extracted from jwt claims")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, int32(userID))
		request := r.WithContext(ctx)

		next.ServeHTTP(w, request)

	})
}

func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)

	})
}
