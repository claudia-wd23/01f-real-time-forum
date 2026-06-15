package router

import (
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/handlers"
	"real-time-forum/backend/middleware"
	"real-time-forum/backend/websocket"

	"github.com/gorilla/mux"
)

func SetupRouter(db *database.Database) *mux.Router {
    r := mux.NewRouter()

    api := r.PathPrefix("/api").Subrouter()

    // Public
    api.HandleFunc("/test", handlers.TestHandler).Methods("GET")
    api.HandleFunc("/ws-test", handlers.WebSocketTest)

    api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Backend is running"))
    }).Methods("GET")

    // Auth
    authHandler := &handlers.AuthHandler{DB: db}
    api.HandleFunc("/register", authHandler.Register).Methods("POST")
    api.HandleFunc("/login", authHandler.Login).Methods("POST")
    api.HandleFunc("/logout", authHandler.Logout).Methods("POST")

    // Authenticated routes
    authMiddleware := &middleware.AuthMiddleware{DB: db}

    api.Handle("/me",
        authMiddleware.RequireAuth(http.HandlerFunc(authHandler.CurrentUser)),
    )

    // Posts
    postsHandler := &handlers.PostsHandler{DB: db}
    api.Handle("/posts/create",
        authMiddleware.RequireAuth(http.HandlerFunc(postsHandler.CreatePost)),
    ).Methods("POST")

    api.HandleFunc("/posts", postsHandler.ListPosts).Methods("GET")
    api.HandleFunc("/posts/get", postsHandler.GetPost).Methods("GET")
    api.HandleFunc("/categories", postsHandler.GetCategories).Methods("GET")

    api.Handle("/posts/update",
        authMiddleware.RequireAuth(http.HandlerFunc(postsHandler.UpdatePost)),
    ).Methods("PUT")

    api.Handle("/posts/delete",
        authMiddleware.RequireAuth(http.HandlerFunc(postsHandler.DeletePost)),
    ).Methods("DELETE")

    // Comments
    commentsHandler := &handlers.CommentsHandler{DB: db}

    api.Handle("/comments/create",
        authMiddleware.RequireAuth(http.HandlerFunc(commentsHandler.CreateComment)),
    ).Methods("POST")

    api.HandleFunc("/comments", commentsHandler.ListComments).Methods("GET")

    api.Handle("/comments/update",
        authMiddleware.RequireAuth(http.HandlerFunc(commentsHandler.UpdateComment)),
    ).Methods("PUT")

    api.Handle("/comments/delete",
        authMiddleware.RequireAuth(http.HandlerFunc(commentsHandler.DeleteComment)),
    ).Methods("DELETE")

    // Messages
    messagesHandler := &handlers.MessagesHandler{DB: db}

    api.Handle("/messages/send",
        authMiddleware.RequireAuth(http.HandlerFunc(messagesHandler.SendMessage)),
    ).Methods("POST")

    api.Handle("/messages",
        authMiddleware.RequireAuth(http.HandlerFunc(messagesHandler.ListConversation)),
    ).Methods("GET")

    api.Handle("/messages/read",
        authMiddleware.RequireAuth(http.HandlerFunc(messagesHandler.MarkRead)),
    ).Methods("POST")

    // WebSocket (must be before the static file catch-all)
    r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        hub := websocket.GlobalHub
        websocket.ServeWs(hub, w, r)
    })

    // Serve frontend (static files) — catch-all, must be last
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))

    return r
}

/*package router        //This replaced for the above which has comments on session's functions

import (
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/handlers"
	"real-time-forum/backend/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter(db *database.Database) *mux.Router {
    r := mux.NewRouter()
         r.HandleFunc("/test", handlers.TestHandler).Methods("GET")
         r.HandleFunc("/ws-test", handlers.WebSocketTest)

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Backend is running"))
}).Methods("GET")

    authHandler := &handlers.AuthHandler{DB: db}
    authMiddleware := &middleware.AuthMiddleware{DB: db}

    r.HandleFunc("/register", authHandler.Register).Methods("POST")
    r.HandleFunc("/login", authHandler.Login).Methods("POST")
    r.HandleFunc("/logout", authHandler.Logout).Methods("POST")

    r.Handle("/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.CurrentUser)))

    postsHandler := &handlers.PostsHandler{DB: db}

    r.Handle("/posts/create",
        authMiddleware.RequireAuth(http.HandlerFunc(postsHandler.CreatePost)),
    ).Methods("POST")

    r.HandleFunc("/posts", postsHandler.ListPosts).Methods("GET")
    r.HandleFunc("/posts/get", postsHandler.GetPost).Methods("GET")

    r.Handle("/posts/update",
        authMiddleware.RequireAuth(http.HandlerFunc(postsHandler.UpdatePost)),
    ).Methods("PUT")

    r.Handle("/posts/delete",
        authMiddleware.RequireAuth(http.HandlerFunc(postsHandler.DeletePost)),
    ).Methods("DELETE")

    commentsHandler := &handlers.CommentsHandler{DB: db}

    r.Handle("/comments/create",
        authMiddleware.RequireAuth(http.HandlerFunc(commentsHandler.CreateComment)),
    ).Methods("POST")

    r.HandleFunc("/comments", commentsHandler.ListComments).Methods("GET")

    r.Handle("/comments/update",
        authMiddleware.RequireAuth(http.HandlerFunc(commentsHandler.UpdateComment)),
    ).Methods("PUT")

    r.Handle("/comments/delete",
        authMiddleware.RequireAuth(http.HandlerFunc(commentsHandler.DeleteComment)),
    ).Methods("DELETE")

    messagesHandler := &handlers.MessagesHandler{DB: db}

    r.Handle("/messages/send",
        authMiddleware.RequireAuth(http.HandlerFunc(messagesHandler.SendMessage)),
    ).Methods("POST")

    r.Handle("/messages",
        authMiddleware.RequireAuth(http.HandlerFunc(messagesHandler.ListConversation)),
    ).Methods("GET")

    r.Handle("/messages/read",
        authMiddleware.RequireAuth(http.HandlerFunc(messagesHandler.MarkRead)),
    ).Methods("POST")

    return r
}*/
