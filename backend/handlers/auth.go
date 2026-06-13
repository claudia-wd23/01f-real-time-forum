package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

type AuthHandler struct {
    DB *database.Database
}

type userResponse struct {
    ID       int    `json:"id"`
    Nickname string `json:"nickname"`
    Email    string `json:"email"`
}

type registerRequest struct {
    Nickname  string `json:"nickname"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Age       int    `json:"age"`
    Gender    string `json:"gender"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if req.Nickname == "" || req.Email == "" || req.Password == "" {
        http.Error(w, "nickname, email and password are required", http.StatusBadRequest)
        return
    }

    user, err := h.DB.CreateUser(req.Nickname, req.Email, req.Password)
    if err != nil {
        http.Error(w, "could not create user", http.StatusBadRequest)
        return
    }

    session, err := h.DB.CreateSession(user.ID)
    if err != nil {
        http.Error(w, "could not create session", http.StatusInternalServerError)
        return
    }

    utils.SetSessionCookie(w, session.ID)

    json.NewEncoder(w).Encode(userResponse{ID: user.ID, Nickname: user.Username, Email: user.Email})
}

type loginRequest struct {
    Identifier string `json:"identifier"`
    Password   string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    user, err := h.DB.AuthenticateUser(req.Identifier, req.Password)
    if err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    session, err := h.DB.CreateSession(user.ID)
    if err != nil {
        http.Error(w, "could not create session", http.StatusInternalServerError)
        return
    }

    utils.SetSessionCookie(w, session.ID)

    json.NewEncoder(w).Encode(userResponse{ID: user.ID, Nickname: user.Username, Email: user.Email})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    sessionID, err := utils.GetSessionIDFromRequest(r)
    if err == nil {
        _ = h.DB.DeleteSession(sessionID)
    }

    utils.ClearSessionCookie(w)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("logged out"))
}

func (h *AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
    raw := r.Context().Value("user")
    if raw == nil {
        http.Error(w, "not authenticated", http.StatusUnauthorized)
        return
    }

    u, ok := raw.(*models.User)
    if !ok {
        http.Error(w, "not authenticated", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(userResponse{ID: u.ID, Nickname: u.Username, Email: u.Email})
}
