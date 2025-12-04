package handlers

import (
	"database/sql"
	"encoding/json"
	"farmequip_api/models"
	"net/http"
)

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req models.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Gagal membaca request body", 400)
			return
		}

		if (req.Email == "" && req.Username == "") || req.Password == "" {
			http.Error(w, "Email/Username dan Password wajib diisi", 400)
			return
		}

		var user models.User
		err := db.QueryRow(`
			SELECT id, nama, email, username 
			FROM users 
			WHERE (email = ? OR username = ?) AND password = ?
		`, req.Email, req.Username, req.Password).
			Scan(&user.ID, &user.Nama, &user.Email, &user.Username)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Email/Username atau password salah", 400)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"user":   user,
		})
	}
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Parameter id diperlukan", 400)
			return
		}

		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Gagal membaca request body", 400)
			return
		}

		_, err := db.Exec(`
			UPDATE users 
			SET nama = ?, email = ?, username = ?, password = ?
			WHERE id = ?
		`, u.Nama, u.Email, u.Username, u.Password, id)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write([]byte("User berhasil diupdate"))
	}
}
