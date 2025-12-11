# FarmEquip API

API untuk manajemen sewa alat pertanian yang dibangun menggunakan Go

## Teknologi yang Digunakan

- **Go** - Bahasa pemrograman
- **Gorilla Mux** - HTTP router dan URL matcher
- **MySQL** - Database
- **Cloudinary** - Cloud storage untuk gambar
- **Functional Programming Principles** 

## Fitur Utama

- Autentikasi pengguna (Login)
- Manajemen kategori alat pertanian
- Manajemen alat pertanian (CRUD)
- Upload gambar ke Cloudinary
- Sorting dan filtering data
- CORS enabled untuk integrasi frontend

## API Endpoints

### Authentication

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

**Response:**
```json
{
  "status": "success",
  "user": {
    "id": 1,
    "nama": "John Doe",
    "email": "user@example.com",
    "username": "username"
  }
}
```

---

### Kategori

#### Get All Kategori
```http
GET /kategori
```

**Response:**
```json
[
  {
    "id": 1,
    "nama_kategori": "Traktor",
    "deskripsi": "Alat untuk membajak sawah",
    "slug": "traktor"
  }
]
```

#### Create Kategori
```http
POST /kategori
Content-Type: application/json

{
  "nama_kategori": "Traktor",
  "deskripsi": "Alat untuk membajak sawah"
}
```

**Response:**
```
Kategori berhasil ditambahkan
```

#### Update Kategori
```http
PUT /kategori?id=1
Content-Type: application/json

{
  "nama_kategori": "Traktor Besar",
  "deskripsi": "Alat untuk membajak sawah luas"
}
```

**Response:**
```
Kategori berhasil diupdate
```

#### Delete Kategori
```http
DELETE /kategori?id=1
```

**Response:**
```
Kategori berhasil dihapus
```

---

### Alat Pertanian

#### Get All Alat
```http
GET /alat
```

**Query Parameters:**
- `sort` (optional): 
  - `nama_asc` - Sort berdasarkan nama A-Z
  - `nama_desc` - Sort berdasarkan nama Z-A
  - `harga_asc` - Sort berdasarkan harga terendah
  - `harga_desc` - Sort berdasarkan harga tertinggi
  - `newest` - Sort berdasarkan terbaru
  - `oldest` - Sort berdasarkan terlama

**Example:**
```http
GET /alat?sort=harga_asc
```

**Response:**
```json
[
  {
    "id": 1,
    "nama_alat": "Traktor Kubota",
    "kategori_id": 1,
    "nama_kategori": "Traktor",
    "deskripsi": "Traktor dengan kapasitas besar",
    "harga_per_hari": 500000,
    "harga_per_minggu": 3000000,
    "harga_per_bulan": 10000000,
    "gambar": "https://cloudinary.com/...",
    "spesifikasi": "Mesin diesel 100HP"
  }
]
```

#### Get Alat by ID
```http
GET /alat/{id}
```

**Example:**
```http
GET /alat/1
```

#### Get Alat by Kategori Slug
```http
GET /alat/{slug}
```

**Example:**
```http
GET /alat/traktor
```

#### Create Alat
```http
POST /alat
Content-Type: multipart/form-data

Form Data:
- nama_alat: "Traktor Kubota"
- kategori_id: 1
- deskripsi: "Traktor dengan kapasitas besar"
- harga_per_hari: 500000
- harga_per_minggu: 3000000
- harga_per_bulan: 10000000
- spesifikasi: "Mesin diesel 100HP"
- gambar: [file upload]
```

**Response:**
```
Alat berhasil ditambahkan
```

#### Update Alat
```http
PUT /alat?id=1
Content-Type: multipart/form-data

Form Data:
- nama_alat: "Traktor Kubota M135"
- kategori_id: 1
- deskripsi: "Traktor dengan kapasitas sangat besar"
- harga_per_hari: 600000
- harga_per_minggu: 3500000
- harga_per_bulan: 12000000
- spesifikasi: "Mesin diesel 135HP"
- gambar: [file upload] (optional)
```

**Response:**
```
Alat berhasil diperbarui
```

#### Delete Alat
```http
DELETE /alat?id=1
```

**Response:**
```
Alat berhasil dihapus
```

---

## Error Responses

API akan mengembalikan HTTP status code yang sesuai:

- `200` - Success
- `400` - Bad Request (validasi gagal)
- `500` - Internal Server Error

## CORS Configuration

API ini sudah dikonfigurasi dengan CORS yang mengizinkan:
- Origin: `*` (semua origin)
- Methods: `GET, POST, PUT, DELETE, OPTIONS`
- Headers: `Content-Type, Authorization`
