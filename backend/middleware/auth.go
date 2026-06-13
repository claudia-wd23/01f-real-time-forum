package middleware

import (
	"context"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
)

type AuthMiddleware struct {
    DB *database.Database
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        sessionID, err := utils.GetSessionIDFromRequest(r)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        session, err := m.DB.GetSession(sessionID)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        user, err := m.DB.GetUserByID(session.UserID)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
