# API Web Service

API Web Service sederhana dibangun dengan **Go (Golang)** untuk menyediakan RESTful endpoints yang digunakan untuk berinteraksi dengan data melalui HTTP.

Struktur folder utama mencakup:
- `config/` – konfigurasi aplikasi
- `handlers/` – fungsi penanganan request HTTP
- `middleware/` – middleware untuk logging/auth
- `models/` – definisi model data
- `routes/` – definisi route API
- `utils/` – utilitas dan helper
- `errors/` – error handling custom

## 🧠 Teknologi Utama

- Go (Golang)
- RESTful API style
- Modular architecture (handlers, models, routes, middleware)
- JSON sebagai format response

## 🚀 Fitur

✔ Endpoint CRUD (sesuai implementasi di `routes/`)  
✔ Middleware untuk validasi dan logging  
✔ Struktur folder terpisah untuk clean architecture

## 📦 Instalasi

1. Clone repository ini  
```bash
git clone https://github.com/lkswilliam2009/web-service.git
cd web-service
   ```
2. Install dependencies
```bash
go mod download
```

## ☕ Menjalankan Server

Jalankan server dengan:
```bash
go run main.go
```

Atau gunakan:

```bash
go build
./web-service
```

## 📡 Contoh Request

Contoh menggunakan curl:

```bash
curl http://localhost:8080/api
```

(Ubah path dan port sesuai implementasi di routes/)

## 📁 Struktur Direktori

```bash
.
├── config/
├── errors/
├── handlers/
├── middleware/
├── models/
├── routes/
├── utils/
├── main.go
└── README.md
```

## 🧪 Testing

Jika ada folder atau script test:
```bash
go test ./...
```

## 📝 Catatan

Sesuaikan environment variable dan konfigurasi database (jika pakai) di folder config/.

Pastikan Go sudah terinstal di mesin kamu (versi minimal disesuaikan dengan project).

## 💻 Kontribusi

Fork repository

Buat branch fitur baru

Lakukan perubahan

Submit Pull Request

## 📄 License

Lisensi tergantung apa yang ditetapkan di project ini (jika belum ada, bisa pakai MIT atau lainnya).