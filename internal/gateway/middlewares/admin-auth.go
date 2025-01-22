package middlewares

import (
	"context"
	"net/http"
)

var adminCtxKey = &contextKey{"admin"}

func AdminAuthMiddleware(adminSecretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			adminKey := r.Header.Get("x-admin-key")
			if adminKey == "" {
				next.ServeHTTP(w, r)
				return
			}
			if adminKey != adminSecretKey {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), adminCtxKey, adminKey)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForAdminContext(ctx context.Context) string {
	raw, ok := ctx.Value(adminCtxKey).(string)
	if !ok {
		return ""
	}
	return raw
}
