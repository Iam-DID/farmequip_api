package main

import (
    "log"
    "net/http"
    "os"
    "farmequip_api/database"
)

func main() {
    db := database.ConnectDB()
    defer db.Close()

    // Init Cloudinary
    cld, err := database.InitCloudinary()
    if err != nil {
        log.Fatal("Cloudinary gagal init:", err)
    }

    // Check Cloudinary connection ---- HERE
    if err := database.CheckCloudinaryConnection(cld); err != nil {
        log.Fatal(err)
    }
    log.Println("Cloudinary connection OK")

    Setupcld(cld)
    SetupRoutes(db, cld)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Println("Server berjalan di port:", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}