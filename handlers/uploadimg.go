package handlers

import (
    "context"
    "fmt"
    "net/http"

    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadHandler struct {
    Cloud *cloudinary.Cloudinary
}

func NewUploadHandler(cld *cloudinary.Cloudinary) *UploadHandler {
    return &UploadHandler{Cloud: cld}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Image not found", http.StatusBadRequest)
        return
    }
    defer file.Close()

    uploadResult, err := h.Cloud.Upload.Upload(
        context.Background(),
        file,
        uploader.UploadParams{
            Folder:   "/",
            PublicID: header.Filename,
        },
    )
    if err != nil {
        http.Error(w, "Failed to upload", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"url": "%s"}`, uploadResult.SecureURL)
}

