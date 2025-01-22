package middlewares

import (
	"context"
	"ecommerce-microservices/pkg/jwt"
	"net/http"
	"strings"
)

var userCtxKey = &contextKey{"user"}

func JWTAuthMiddleware(jwtSecretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(header) != 2 {
				next.ServeHTTP(w, r)
				return
			}

			token := header[1]
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			username, err := jwt.ParseToken(token, jwtSecretKey)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, username)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForJWTContext(ctx context.Context) string {
	raw, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		return ""
	}
	return raw
}
