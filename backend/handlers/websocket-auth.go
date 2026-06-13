package handlers

import (
	"net/http"
	"real-time-forum/backend/websocket"
	"strconv"
	//"real-time-forum/backend/handlers"
)

func WebSocketAuth(hub *websocket.Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userIDStr := r.Header.Get("X-User-ID")
        userID, _ := strconv.Atoi(userIDStr)

        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            return
        }

        client := &websocket.Client{
            UserID: userID,
            Conn:   conn,
            Send:   make(chan []byte),
        }

        hub.Register <- client

        go client.ReadPump()
        go client.WritePump()
    }
}
