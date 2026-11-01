#  API Beware of Shark Attack

API Serverless berbasis **Go (Golang)** yang terintegrasi dengan database **PostgreSQL (Neon Cloud DB)** dan dikonfigurasi untuk deployment di platform **Vercel**. API ini menyediakan manajemen data histori insiden serangan hiu di berbagai belahan dunia, lokasi kejadian, spesies hiu, dan aktivitas korban.

- **URL Deployment Production**: `https://shark-attack-api-git-main-z0rnos.vercel.app`

---

##  Kegunaan & Fitur Utama

- **Akses Data Insiden**: Mengambil daftar insiden serangan hiu beserta relasi data lokasi, spesies hiu, dan aktivitas dengan dukungan paginasi (`limit` dan `page`).
- **Autentikasi Pengguna**: Fitur Registrasi & Login terenkripsi menggunakan `bcrypt` untuk enkripsi password dan penandatanganan **JSON Web Token (JWT)**.
- **Endpoint Terproteksi (JWT)**: Pengamanan area profil user dan manajemen data aktivitas (*CRUD*) menggunakan middleware JWT Bearer.
- **Admin & Database Seeder (Basic Auth)**: Endpoint admin terproteksi **Basic Authentication** untuk menjalankan proses *import/seeding* data awal dari file CSV (`data.csv`) ke PostgreSQL Neon.
- **Vercel Serverless Function**: Arsitektur yang disesuaikan untuk berjalan di lingkungan Vercel Serverless (`api/index.go`) maupun server lokal (`main.go`).

---

##  Persyaratan Sistem & Dependensi

- **Go**: Version 1.20+
- **Database**: PostgreSQL (Neon Serverless PostgreSQL)
- **Deployment Platform**: Vercel
- **Dependencies Utama**:
  - `github.com/lib/pq` (Driver PostgreSQL)
  - `github.com/golang-jwt/jwt/v5` (Manajemen JWT)
  - `golang.org/x/crypto` (Password Hashing via bcrypt)

---

##  Konfigurasi Environment Variables (`.env`)

Buat file `.env` di direktori utama project atau atur pada **Vercel Project Settings > Environment Variables**:

| Variable | Deskripsi | Contoh Value |
| :--- | :--- | :--- |
| `DATABASE_URL` | Connection String PostgreSQL (Neon DB) | `postgresql://user:pass@ep-xyz.aws.neon.tech/neondb?sslmode=require` |
| `POSTGRES_URL` | Fallback URL bawaan dari integrasi Vercel | `postgresql://...` |
| `JWT_SECRET` | Secret key untuk enkripsi token JWT | `super_secret_jwt_key_shark_attack` |
| `BASIC_AUTH_USER` | Username untuk akses endpoint Admin | `admin` |
| `BASIC_AUTH_PASS` | Password untuk akses endpoint Admin | `admin123` |
| `PORT` | Port server saat dijalankan di lokal | `8081` |

---

##  Skema Tabel Database (PostgreSQL)

Jalankan perintah SQL DDL berikut di Dashboard Neon SQL Editor sebelum melakukan seeding data:

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    country VARCHAR(100),
    area VARCHAR(150),
    location_name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS shark_species (
    id SERIAL PRIMARY KEY,
    species_name VARCHAR(150),
    species_group VARCHAR(100) DEFAULT 'General'
);

CREATE TABLE IF NOT EXISTS activities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150),
    category VARCHAR(100) DEFAULT 'Water Activity',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    incident_date VARCHAR(100),
    year INT,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    location_id INT REFERENCES locations(id) ON DELETE SET NULL,
    species_id INT REFERENCES shark_species(id) ON DELETE SET NULL,
    activity_id INT REFERENCES activities(id) ON DELETE SET NULL,
    injury TEXT,
    fatal VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

Panduan Penggunaan & Running1. Menjalankan Server LokalBash# Clone repository dan masuk ke folder project
cd beware-Of-Shark-Attack

# Install modul dependencies
go mod download

# Menjalankan server lokal
go run main.go
Server akan aktif di http://localhost:8081.2. Menjalankan Seeding Data CSVLakukan request HTTP POST ke endpoint berikut (menggunakan Postman/Thunder Client):URL: https://shark-attack-api-git-main-z0rnos.vercel.app/api/admin/seed-csv (Production) atau http://localhost:8081/api/admin/seed-csv (Lokal)Auth: Select Basic AuthUsername: adminPassword: admin123🛣️ Daftar Path & Endpoint API🌐 1. Root / General EndpointGET /Akses: PublicDeskripsi: Mengembalikan status API dan 10 data insiden terbaru dari database Neon PostgreSQL.Response Success (200 OK):JSON{
  "status": "success",
  "total": 10,
  "data": [
    {
      "id": 1,
      "incident_date": "11/06/2017",
      "year": 2017,
      "location_id": 1,
      "species_id": 1,
      "activity_id": 1,
      "injury": "No injury, board bitten",
      "fatal": "N",
      "created_at": "2026-03-30T10:00:00Z",
      "location": {
        "id": 1,
        "country": "AUSTRALIA",
        "area": "Western Australia",
        "location_name": "Point Casuarina, Bunbury"
      },
      "species": {
        "id": 1,
        "species_name": "White shark, 4 m",
        "species_group": "General"
      },
      "activity": {
        "id": 1,
        "name": "Bodyboarding",
        "category": "Water Activity"
      }
    }
  ]
}
2. User Authentication EndpointsPOST /api/users/registerAkses: PublicRequest Body:JSON{
  "username": "asepAkamsi",
  "email": "asep@akamsi.com",
  "password": "123"
}
Response Success (201 Created):JSON{
  "status": "success",
  "message": "Registrasi berhasil",
  "data": {
    "id": 1,
    "username": "asepAkamsi",
    "email": "asep@Akamsi.com",
    "role": "user",
    "created_at": "2026-03-30T10:00:00Z"
  }
}
POST /api/users/loginAkses: PublicRequest Body:JSON{
  "email": "asep@Akamsi.com",
  "password": "123"
}
Response Success (200 OK):JSON{
  "status": "success",
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
3. Incidents EndpointsGET /api/incidentsAkses: PublicQuery Parameters:limit (optional): Jumlah data per halaman (Default: 10)page (optional): Halaman ke-n (Default: 1)Response Success (200 OK): Menampilkan list objek insiden.GET /api/incidents/{id}Akses: PublicResponse Success (200 OK): Menampilkan detail 1 data insiden berdasarkan ID.POST /api/incidentsAkses: Public / JWT OptionalRequest Body:JSON{
  "incident_date": "15/08/2023",
  "year": 2023,
  "location_id": 1,
  "species_id": 1,
  "activity_id": 1,
  "injury": "Minor laceration",
  "fatal": "N"
}
DELETE /api/incidents/{id}Akses: Public / JWT OptionalResponse Success (200 OK): {"message": "Berhasil menghapus data insiden"}🏊 4. Activities Endpoints (Terproteksi JWT)Catatan: Semua endpoint /api/activities membutuhkan Header Authorization: Bearer <token>.MethodPathDeskripsiGET/api/activitiesMengambil daftar seluruh aktivitasGET/api/activities/{id}Mengambil detail aktivitas berdasarkan IDPOST/api/activitiesMenambahkan data aktivitas baruPUT/api/activities/{id}Memperbarui data aktivitas berdasarkan IDDELETE/api/activities/{id}Menghapus data aktivitasRequest Body untuk POST/PUT:JSON{
  "name": "Kite Surfing",
  "category": "Water Sports"
}
5. User Profile EndpointGET /api/profileAkses: Protected JWT (Authorization: Bearer <token>)Response Success (200 OK):JSON{
  "message": "Akses area terproteksi JWT berhasil",
  "user": {
    "user_id": 1,
    "username": "johndoe",
    "role": "user"
  }
}
6. Admin EndpointsGET /api/admin/dashboardAkses: Basic Auth (BASIC_AUTH_USER & BASIC_AUTH_PASS)Response Success (200 OK): {"message": "Akses admin terproteksi Basic Auth berhasil"}POST /api/admin/seed-csvAkses: Basic AuthResponse Success (200 OK):JSON{
  "status": "success",
  "message": "Import CSV data insiden berhasil",
  "total": 50
}
