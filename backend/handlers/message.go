package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
)

type MessagesHandler struct {
    DB *database.Database
}

type sendMessageRequest struct {
    ReceiverID int    `json:"receiver_id"`
    Content    string `json:"content"`
}

func (h *MessagesHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    var req sendMessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    msg, err := h.DB.SendMessage(user.ID, req.ReceiverID, req.Content)
    if err != nil {
        http.Error(w, "could not send message", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(msg)
}

func (h *MessagesHandler) ListConversation(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)

    otherIDStr := r.URL.Query().Get("user_id")
    otherID, _ := strconv.Atoi(otherIDStr)

    msgs, err := h.DB.ListMessagesBetween(user.ID, otherID)
    if err != nil {
        http.Error(w, "could not load messages", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(msgs)
}

func (h *MessagesHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    if err := h.DB.MarkMessageRead(id); err != nil {
        http.Error(w, "message not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
}
