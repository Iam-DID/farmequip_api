package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"farmequip_api/models"
	"net/http"
	"sort"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
)

func GetAlat(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        sortParam := r.URL.Query().Get("sort")

        rows, err := db.Query(`
            SELECT 
                a.id,
                a.nama_alat,
                a.kategori_id,
                k.nama_kategori,
                a.deskripsi,
                a.harga_per_hari,
                a.harga_per_minggu,
                a.harga_per_bulan,
                a.gambar,
                a.spesifikasi
            FROM alat_pertanian a
            JOIN kategori k ON k.id = a.kategori_id
        `)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        defer rows.Close()

        var list []models.Alat

        for rows.Next() {
            var a models.Alat
            err := rows.Scan(
                &a.ID,
                &a.NamaAlat,
                &a.KategoriID,
                &a.NamaKategori,
                &a.Deskripsi,
                &a.HargaHarian,
                &a.HargaMingguan,
                &a.HargaBulanan,
                &a.Gambar, // sekarang string URL
                &a.Spesifikasi,
            )
            if err != nil {
                http.Error(w, err.Error(), 500)
                return
            }

            list = append(list, a)
        }

        sort.Slice(list, func(i, j int) bool {
            switch sortParam {
            case "nama_asc":
                return list[i].NamaAlat < list[j].NamaAlat
            case "nama_desc":
                return list[i].NamaAlat > list[j].NamaAlat
            case "harga_asc":
                return list[i].HargaHarian < list[j].HargaHarian
            case "harga_desc":
                return list[i].HargaHarian > list[j].HargaHarian
            case "newest":
                return list[i].ID > list[j].ID
            case "oldest":
                return list[i].ID < list[j].ID
            default:
                return list[i].ID < list[j].ID
            }
        })

        json.NewEncoder(w).Encode(list)
    }
}


func GetToolById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := mux.Vars(r)["id"]

		rows, err := db.Query(`
            SELECT 
                a.id,
                a.nama_alat,
                a.kategori_id,
				k.nama_kategori,
                a.deskripsi,
                a.harga_per_hari,
                a.harga_per_minggu,
                a.harga_per_bulan,
                a.gambar,
				a.spesifikasi
            FROM alat_pertanian a
            JOIN kategori k ON a.kategori_id = k.id
            WHERE a.id = ?
        `, id)

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()

		var list []models.Alat

		for rows.Next() {
			var a models.Alat
			var imgBytes []byte

			err := rows.Scan(
				&a.ID,
				&a.NamaAlat,
				&a.KategoriID,
				&a.NamaKategori,
				&a.Deskripsi,
				&a.HargaHarian,
				&a.HargaMingguan,
				&a.HargaBulanan,
				&imgBytes,
				&a.Spesifikasi,
			)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			a.Gambar = base64.StdEncoding.EncodeToString(imgBytes)

			list = append(list, a)
		}

		if len(list) == 0 {
			w.Write([]byte("Tidak ada alat dalam kategori ini"))
			return
		}

		json.NewEncoder(w).Encode(list)
	}
}

func CreateAlat(db *sql.DB, cld *cloudinary.Cloudinary) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        r.ParseMultipartForm(10 << 20)

        nama := r.FormValue("nama_alat")
        kategori := r.FormValue("kategori_id")
        deskripsi := r.FormValue("deskripsi")
        hari := r.FormValue("harga_per_hari")
        minggu := r.FormValue("harga_per_minggu")
        bulan := r.FormValue("harga_per_bulan")
        spesifikasi := r.FormValue("spesifikasi")

        file, header, err := r.FormFile("gambar")
        if err != nil {
            http.Error(w, "Gambar wajib diupload", 400)
            return
        }
        defer file.Close()

        // Upload Cloudinary
        upload, err := cld.Upload.Upload(
            r.Context(),
            file,
            uploader.UploadParams{
                Folder:   "/",
                PublicID: header.Filename,
            },
        )
        if err != nil {
            http.Error(w, "Upload Cloudinary gagal: "+err.Error(), 500)
            return
        }

        _, err = db.Exec(`
            INSERT INTO alat_pertanian 
            (nama_alat, kategori_id, deskripsi, harga_per_hari, harga_per_minggu, harga_per_bulan, gambar, spesifikasi)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        `, nama, kategori, deskripsi, hari, minggu, bulan, upload.SecureURL, spesifikasi)

        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        w.Write([]byte("Alat berhasil ditambahkan"))
    }
}


func GetAlatBySlug(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		slug := mux.Vars(r)["slug"]
		if slug == "" {
			w.Write([]byte("Slug kategori wajib diisi"))
			return
		}

		rows, err := db.Query(`
            SELECT 
                a.id,
                a.nama_alat,
                a.kategori_id,
                k.nama_kategori,
                a.deskripsi,
                a.harga_per_hari,
                a.harga_per_minggu,
                a.harga_per_bulan,
                a.gambar,
				a.spesifikasi
            FROM alat_pertanian a
            JOIN kategori k ON a.kategori_id = k.id
            WHERE k.slug = ?
        `, slug)

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()

		var list []models.Alat

		for rows.Next() {
			var a models.Alat
			var imgBytes []byte

			err := rows.Scan(
				&a.ID,
				&a.NamaAlat,
				&a.KategoriID,
				&a.NamaKategori,
				&a.Deskripsi,
				&a.HargaHarian,
				&a.HargaMingguan,
				&a.HargaBulanan,
				&imgBytes,
				&a.Spesifikasi,
			)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			a.Gambar = base64.StdEncoding.EncodeToString(imgBytes)

			list = append(list, a)
		}

		if len(list) == 0 {
			w.Write([]byte("Tidak ada alat dalam kategori ini"))
			return
		}

		json.NewEncoder(w).Encode(list)
	}
}

func UpdateAlat(db *sql.DB, cld *cloudinary.Cloudinary) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        id := r.URL.Query().Get("id")
        if id == "" {
            http.Error(w, "ID alat wajib diisi", 400)
            return
        }

        r.ParseMultipartForm(10 << 20)

        nama := r.FormValue("nama_alat")
        kategori := r.FormValue("kategori_id")
        deskripsi := r.FormValue("deskripsi")
        hari := r.FormValue("harga_per_hari")
        minggu := r.FormValue("harga_per_minggu")
        bulan := r.FormValue("harga_per_bulan")
        spesifikasi := r.FormValue("spesifikasi")

        // Ambil gambar lama
        var oldURL string
        db.QueryRow(`SELECT gambar FROM alat_pertanian WHERE id = ?`, id).Scan(&oldURL)

        // Cek ada upload baru?
        file, header, err := r.FormFile("gambar")

        newURL := oldURL

        if err == nil {
            // Ada gambar baru → upload Cloudinary
            defer file.Close()

            upload, err := cld.Upload.Upload(
                r.Context(),
                file,
                uploader.UploadParams{
                    Folder:   "/",
                    PublicID: header.Filename,
                },
            )
            if err != nil {
                http.Error(w, "Gagal upload cloudinary: "+err.Error(), 500)
                return
            }
            newURL = upload.SecureURL
        }

        _, err = db.Exec(`
            UPDATE alat_pertanian
            SET nama_alat=?, kategori_id=?, deskripsi=?,
                harga_per_hari=?, harga_per_minggu=?, harga_per_bulan=?,
                gambar=?, spesifikasi=?
            WHERE id=?
        `, nama, kategori, deskripsi, hari, minggu, bulan, newURL, spesifikasi, id)

        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        w.Write([]byte("Alat berhasil diperbarui"))
    }
}



func DeleteAlat(db *sql.DB, cld *cloudinary.Cloudinary) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        id := r.URL.Query().Get("id")
        if id == "" {
            http.Error(w, "ID wajib diisi", 400)
            return
        }

        var imageURL string
        db.QueryRow(`SELECT gambar FROM alat_pertanian WHERE id = ?`, id).Scan(&imageURL)

        // hapus cloudinary (jika ada)
        if imageURL != "" {
            // contoh: https://res.cloudinary.com/<cloud>/image/upload/farmequip/nama.jpg
            parts := strings.Split(imageURL, "/")
            publicID := strings.TrimSuffix(parts[len(parts)-1], ".jpg")

            cld.Upload.Destroy(r.Context(), uploader.DestroyParams{
                PublicID: "/" + publicID,
            })
        }

        _, err := db.Exec(`DELETE FROM alat_pertanian WHERE id = ?`, id)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        w.Write([]byte("Alat berhasil dihapus"))
    }
}


