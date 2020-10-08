package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/services/jwt"
)

// EnrichUserID adds the user id to the context
func (middleware Client) EnrichUserID(jwtService jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			//
			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenString := header

			userID, err := jwtService.GetUserIDFromToken(tokenString)
			if err != nil {
				errorResponse, err := json.Marshal(gqlerror.Errorf("invalid token"))
				if err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}

				http.Error(w, string(errorResponse), http.StatusUnauthorized)
				return
			}

			// put it in context
			ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
			ctx = context.WithValue(ctx, ContextKeyJWTToken, tokenString)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
