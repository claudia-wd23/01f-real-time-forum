package main

import (
	"log"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/router"
	"real-time-forum/backend/websocket"
)

	func main() {
    // 1. Open database (includes migrations)
    db, err := database.Open("./forum.db", "./backend/database/migrations.sql")
    if err != nil {
        log.Fatal("Database initialization failed:", err)
    }
// Start WebSocket hub
    go websocket.GlobalHub.Run()

    // 2. Create router
    r := router.SetupRouter(db)

    // 3. Start server
    log.Println("Server running on http://localhost:8080")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}

/*http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    websocket.ServeWs(hub, w, r)
})*/


/*package main

import (
	"log"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/router"
)

func main() {
    // 1. Connect to the database
    err := database.Connect()
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }

    // 2. Run migrations (optional but recommended)
    err = database.RunMigrations()
    if err != nil {
        log.Fatal("Database migrations failed:", err)
    }

    // 3. Create the router
    r := router.SetupRouter()

    // 4. Start the server
    log.Println("Server running on http://localhost:8080")
    err = http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}*/
