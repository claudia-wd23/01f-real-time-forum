package websocket

type Message struct {
    Type      string `json:"type"`       // e.g. "chat"
    SenderID  int    `json:"sender_id"`
    Content   string `json:"content"`
}

