# Backend E-Wallet

![Go](https://img.shields.io/badge/Go-1.26.2-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12.0-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)

REST API backend untuk aplikasi dompet digital (e-wallet). Dibangun dengan Go dan Gin framework, mendukung autentikasi JWT, manajemen transaksi, dan dokumentasi Swagger.

---

## Daftar Isi

- [Fitur](#fitur)
- [Tech Stack](#tech-stack)
- [Struktur Project](#struktur-project)
- [Prasyarat](#prasyarat)
- [Setup Lokal (Tanpa Docker)](#setup-lokal-tanpa-docker)
- [Setup dengan Docker Compose](#setup-dengan-docker-compose)
- [Konfigurasi Environment](#konfigurasi-environment)
- [Konfigurasi Redis](#konfigurasi-redis)
- [API Endpoints](#api-endpoints)
- [Skema Database](#skema-database)
- [Kontribusi](#kontribusi)

---

## Fitur

- Registrasi, login, logout, lupa password, dan reset password via OTP email
- Autentikasi berbasis JWT dengan validasi token aktif (Redis blacklist)
- Manajemen profil pengguna — nama, nomor HP, foto, password, dan PIN
- Dashboard: saldo wallet, total income, total expense
- Top up saldo dengan pemilihan metode pembayaran
- Transfer antar pengguna dengan verifikasi PIN 6 digit
- Riwayat transaksi dengan pencarian dan paginasi
- Laporan grafik transaksi berdasarkan rentang tanggal
- Pencarian penerima transfer berdasarkan nama atau nomor HP
- Serving file statis untuk foto profil
- OTP reset password dengan cooldown 60 detik dan TTL 5 menit
- Dokumentasi Swagger di `/swagger/index.html`

---

## Tech Stack

| Layer | Teknologi |
|---|---|
| Language | Go 1.26.2 |
| Web Framework | Gin 1.12.0 |
| Database | PostgreSQL 17 |
| Cache & Session | Redis 7 (go-redis v9) |
| Auth | JWT (golang-jwt/jwt v5) |
| Password Hashing | Argon2id |
| Database Driver | pgx/v5 (pgxpool) |
| Migration | golang-migrate |
| Email | Gomail v2 (SMTP) |
| API Docs | Swagger (swaggo/swag) |
| Containerization | Docker & Docker Compose |

---

## Struktur Project

```
backend-ewallet/
├── cmd/
│   └── main.go                  # Entry point aplikasi
├── db/
│   ├── migrations/              # File SQL migration (up & down)
│   └── seeders/
│       └── main.go              # Seeder data dummy
├── docs/                        # Generated Swagger docs
├── internal/
│   ├── apperrors/               # Error helper aplikasi
│   ├── binder/                  # Helper binding & validasi request
│   ├── config/
│   │   ├── psql.config.go       # Koneksi PostgreSQL (pgxpool)
│   │   ├── redis.config.go      # Koneksi Redis
│   │   └── gomail.config.go     # Konfigurasi SMTP
│   ├── controller/              # Handler HTTP (auth, user, transaction)
│   ├── dto/                     # Request & Response DTO
│   ├── jwttoken/                # JWT utility
│   ├── middleware/
│   │   ├── auth.middleware.go   # Validasi JWT + cek blacklist Redis
│   │   └── cors.middleware.go   # Konfigurasi CORS
│   ├── model/                   # Domain model (entity)
│   ├── repository/              # Akses database & Redis
│   ├── response/                # Helper JSON response standar
│   ├── router/                  # Definisi route
│   └── service/                 # Business logic
├── pkg/
│   ├── hash.pkg.go              # Argon2id hash & verify
│   ├── jwt.pkg.go               # Generate & verify JWT
│   ├── gomail.pkg.go            # Kirim email OTP
│   ├── generateOTP.pkg.go       # Generator OTP kriptografis
│   └── validation.pkg.go        # Validasi format email
├── public/
│   └── img/profiles/            # Upload foto profil
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

---

## Prasyarat

Pastikan tools berikut sudah tersedia sebelum menjalankan project:

- Go 1.26.2 atau lebih baru
- PostgreSQL
- Redis
- `golang-migrate` (untuk migrasi manual)
- Make
- Git
- Docker & Docker Compose (opsional, untuk setup berbasis container)

---

## Setup Lokal (Tanpa Docker)

### 1. Clone repository

```bash
git clone https://github.com/aqilknz/backend-ewallet.git
cd backend-ewallet
```

### 2. Install dependensi Go

```bash
go mod download
```

### 3. Buat file `.env`

Buat file `.env` di root project:

```env
# Server
APP_HOST=localhost
APP_PORT=9000

# Database
POSTGRES_HOST=localhost
POSTGRES_USER=myuser
POSTGRES_PASSWORD=12345
POSTGRES_DB=db_ewallet
POSTGRES_PORT=5432

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DB=0

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM_EMAIL=no-reply@ewallet.com

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_ISSUER=aqilknz
```

### 4. Jalankan migrasi database

```bash
make migrate-up
```

### 5. Jalankan seeder (opsional)

```bash
go run ./db/seeders/main.go
```

Seeder akan membuat:
- 5 metode pembayaran (BRI, DANA, BCA, GoPay, OVO)
- 10 user dummy (`user1@mail.com` s.d. `user10@mail.com`, password: `pass1234`, PIN: `123456`)
- 20 transaksi dummy (8 topup, 12 transfer)

### 6. Jalankan server

```bash
make run
```

Server berjalan di:
```
http://localhost:9000
```

Swagger docs:
```
http://localhost:9000/swagger/index.html
```

---

## Setup dengan Docker Compose

Struktur direktori yang diperlukan (`.env`, `redis.conf`, dan `docker-compose.yml` sejajar dengan folder backend dan frontend):

```
project-root/
├── .env
├── redis.conf
├── docker-compose.yml
├── backend-ewallet/        # Folder project ini
└── ewallet-app-react/      # Folder frontend
```

### Jalankan semua service

```bash
docker compose up -d --build
```

Urutan startup yang diatur secara otomatis:
1. `db` — PostgreSQL siap menerima koneksi
2. `redis` — Redis dengan konfigurasi custom `redis.conf`
3. `migrate` — Migrasi database dijalankan sekali lalu selesai
4. `seeder` — Seeder Go dijalankan setelah migrasi berhasil
5. `backend` — Server Go berjalan setelah seeder selesai
6. `web-server` — Frontend Nginx berjalan setelah backend siap

### Port yang digunakan

| Service | Port Host | Port Container |
|---|---|---|
| Backend API | `9000` | `9000` |
| PostgreSQL | `5050` | `5432` |
| Redis | — | `6379` (internal) |
| Frontend | `200` | `200` |

> Redis tidak di-expose ke host secara default. Akses Redis hanya dari dalam jaringan Docker (`ewallet-net`).

### Hentikan semua container

```bash
docker compose down
```

Hapus juga volume (data akan hilang):

```bash
docker compose down -v
```

---

## Konfigurasi Environment

| Variable | Deskripsi | Contoh |
|---|---|---|
| `APP_HOST` | Host server | `0.0.0.0` |
| `APP_PORT` | Port server | `9000` |
| `POSTGRES_HOST` | Host database | `db` (Docker) / `localhost` |
| `POSTGRES_USER` | Username PostgreSQL | `myuser` |
| `POSTGRES_PASSWORD` | Password PostgreSQL | `12345` |
| `POSTGRES_DB` | Nama database | `db_ewallet` |
| `POSTGRES_PORT` | Port PostgreSQL | `5432` |
| `REDIS_HOST` | Host Redis | `redis` (Docker) / `localhost` |
| `REDIS_PORT` | Port Redis | `6379` |
| `REDIS_USERNAME` | Username Redis | `superuser` |
| `REDIS_PASSWORD` | Password Redis | `admin` |
| `REDIS_DB` | Nomor database Redis | `0` |
| `SMTP_HOST` | Host SMTP | `smtp.gmail.com` |
| `SMTP_PORT` | Port SMTP | `587` |
| `SMTP_USER` | Username SMTP | `your_email@gmail.com` |
| `SMTP_PASSWORD` | Password / App Password SMTP | `xxxx xxxx xxxx xxxx` |
| `SMTP_FROM_EMAIL` | Alamat pengirim email | `no-reply@ewallet.com` |
| `JWT_SECRET` | Secret key JWT | `super_secret_key` |
| `JWT_ISSUER` | Issuer JWT | `aqilknz` |
| `VITE_API_URL` | Base URL API untuk frontend | `http://localhost:9000/ewallet` |

> Untuk Gmail, gunakan **App Password** — bukan password akun. Aktifkan 2FA terlebih dahulu di akun Google, lalu buat App Password di [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords).

---

## Konfigurasi Redis

File `redis.conf` diletakkan sejajar dengan `docker-compose.yml`. Berikut konfigurasi yang digunakan:

```conf
# Persistensi RDB
save 3600 1      # setiap jam jika ada 1 perubahan
save 300 100     # setiap 5 menit jika ada 100 perubahan
save 60 10000    # setiap menit jika ada 10000 perubahan

# Persistensi AOF
appendonly yes
appendfsync everysec

# Memory
maxmemory 100mb
maxmemory-policy allkeys-lru

# User (ACL)
user default off
user superuser on >admin ~* &* +@all
```

User `superuser` digunakan oleh aplikasi backend. User `default` dinonaktifkan untuk keamanan.

---

## API Endpoints

Base URL: `/ewallet`

Swagger UI tersedia di `/swagger/index.html`.

### Auth

| Method | Endpoint | Deskripsi | Auth |
|---|---|---|---|
| `POST` | `/auth/register` | Registrasi user baru | — |
| `POST` | `/auth` | Login dan dapatkan JWT token | — |
| `DELETE` | `/auth/logout` | Logout dan invalidasi token | ✅ |
| `POST` | `/auth/create-pin` | Buat PIN 6 digit pertama kali | ✅ |
| `POST` | `/auth/check-email` | Cek apakah email terdaftar | — |
| `POST` | `/auth/forgot-password` | Kirim OTP reset password ke email | — |
| `POST` | `/auth/verify-otp` | Verifikasi kode OTP | — |
| `POST` | `/auth/reset-password` | Reset password setelah OTP valid | — |
| `POST` | `/auth/update-password` | Update password (flow lama) | — |

### Users

| Method | Endpoint | Deskripsi | Auth |
|---|---|---|---|
| `GET` | `/users/profile` | Ambil profil user yang login | ✅ |
| `PATCH` | `/users/profile` | Edit nama, HP, foto profil | ✅ |
| `PATCH` | `/users/profile/password` | Ganti password | ✅ |
| `PATCH` | `/users/profile/pin` | Ganti PIN transaksi | ✅ |
| `GET` | `/users/dashboard` | Ambil saldo, income, expense | ✅ |
| `GET` | `/users/receivers` | Cari penerima transfer | ✅ |

### Transactions

| Method | Endpoint | Deskripsi | Auth |
|---|---|---|---|
| `POST` | `/users/transaction/topup` | Top up saldo | ✅ |
| `POST` | `/users/transaction/transfer` | Transfer ke user lain (butuh PIN) | ✅ |
| `POST` | `/users/transaction/checkpin` | Verifikasi PIN sebelum transaksi | ✅ |
| `GET` | `/users/transaction/history` | Riwayat transaksi dengan paginasi | ✅ |
| `GET` | `/users/transaction/report` | Laporan grafik income/expense | ✅ |

Semua endpoint berproteksi (✅) memerlukan header:
```
Authorization: Bearer <token>
```

---

## Skema Database

```
users
  id, email, password (argon2id), pin (argon2id), created_at, updated_at

profiles
  id, user_id → users, full_name, phone, photo, created_at, updated_at

wallets
  id, user_id → users, balance (BIGINT, ≥ 0), created_at, updated_at

payment_methods
  id, name

transactions
  id, user_id → users, amount, type (topup | transfer_in | transfer_out),
  status (pending | success | failed), created_at, updated_at

topup_details
  id, transaction_id → transactions, payment_method_id → payment_methods,
  discount, tax, sub_total

transfer_details
  id, transaction_id → transactions, receiver_id → users, notes
```

Migrasi dikelola oleh `golang-migrate` dengan file sequential di `db/migrations/`.

### Perintah Migrasi

```bash
# Jalankan semua migrasi
make migrate-up

# Rollback semua migrasi
make migrate-down

# Buat file migrasi baru
make migrate-create NAME=nama_tabel

# Force versi tertentu (jika dirty)
make migrate-force VERSION=1
```

---

## Kontribusi

1. Fork repository ini
2. Buat branch baru dari `master`
3. Buat perubahan dengan pesan commit yang jelas
4. Jalankan project secara lokal dan pastikan endpoint yang terpengaruh berjalan normal
5. Buat pull request dengan deskripsi singkat tentang perubahan yang dilakukan

---

## Related Project

- [Frontend E-Wallet (React)](https://github.com/aqilknz/ewallet-app-react)

---

## License

Project ini dilisensikan di bawah [MIT License](./LICENSE).