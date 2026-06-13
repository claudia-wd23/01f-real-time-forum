package handlers

import (
	"fmt"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Test endpoint working!")
}

/*package handlers

import (
	"encoding/json"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Backend is connected!",
    })
}*/

