package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow all origins for now (you can tighten this later)
        return true
    },
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP → WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    // Create a new client
    client := &Client{
        Hub:    hub,
        Conn:   conn,
        Send:   make(chan []byte, 256),
        UserID: 1, // TEMP: replace with real user ID later
    }

    // Register client with the hub
    hub.Register <- client

    // Start pumps
    go client.writePump()
    go client.readPump()
}
