package websocket

type Hub struct {
    Clients    map[int]*Client
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
}

// GlobalHub is the single hub instance used by the whole app.
var GlobalHub = NewHub()

func NewHub() *Hub {
    return &Hub{
        Clients:    make(map[int]*Client),
        Broadcast:  make(chan []byte),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.Clients[client.UserID] = client

        case client := <-h.Unregister:
            delete(h.Clients, client.UserID)
            close(client.Send)

        case message := <-h.Broadcast:
            for _, client := range h.Clients {
                client.Send <- message
            }
        }
    }

}

/*type Client struct {
    UserID int
    Conn   *websocket.Conn
    Send   chan []byte
}*/

/*package ws

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.Send)
            }

        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.Send <- message:
                default:
                    delete(h.clients, client)
                    close(client.Send)
                }
            }
        }
    }
}*/

/*package ws

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.Send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.Send <- message:
                default:
                    delete(h.clients, client)
                    close(client.Send)
                }
            }
        }
    }
}*/