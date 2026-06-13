package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func WebSocketTest(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("WebSocket upgrade error:", err)
        return
    }
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("WebSocket read error:", err)
            break
        }

        fmt.Println("Received:", string(msg))

        conn.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
    }
}

