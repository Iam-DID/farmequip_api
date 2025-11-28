package main

import (
    "database/sql"
    "farmequip_api/handlers"
    "net/http"

    "github.com/gorilla/mux"
)

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func SetupRoutes(db *sql.DB) {
    r := mux.NewRouter()

    r.Use(CorsMiddleware)

    r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.WriteHeader(http.StatusOK)
    })

    // Login
    r.HandleFunc("/login", handlers.Login(db)).Methods("POST")

    // Users
    r.HandleFunc("/users", handlers.UpdateUser(db)).Methods("PUT")

    // Kategori
    r.HandleFunc("/kategori", handlers.GetKategori(db)).Methods("GET")
    r.HandleFunc("/kategori", handlers.CreateKategori(db)).Methods("POST")
    r.HandleFunc("/kategori", handlers.UpdateKategori(db)).Methods("PUT")
    r.HandleFunc("/kategori", handlers.DeleteKategori(db)).Methods("DELETE")

    // Alat
    r.HandleFunc("/alat", handlers.GetAlat(db)).Methods("GET")
    r.HandleFunc("/alat", handlers.CreateAlat(db)).Methods("POST")
    r.HandleFunc("/alat", handlers.UpdateAlat(db)).Methods("PUT")
    r.HandleFunc("/alat", handlers.DeleteAlat(db)).Methods("DELETE")

    // Detail alat
    r.HandleFunc("/alat/{id:[0-9]+}", handlers.GetToolById(db)).Methods("GET")
    r.HandleFunc("/alat/{slug:[a-zA-Z0-9-_]+}", handlers.GetAlatBySlug(db)).Methods("GET")

    http.Handle("/", r)
}
