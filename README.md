<div align="center">
    <p>
        <a href="https://github.com/404NotFoundIndonesia/" target="_blank">
            <img src="https://avatars.githubusercontent.com/u/87377917?s=200&v=4" width="200" alt="404NFID Logo">
        </a>
    </p>

 [![GitHub Stars](https://img.shields.io/github/stars/iqbaleff214/kamus-banjar-api.svg)](https://github.com/iqbaleff214/kamus-banjar-api/stargazers)
 [![GitHub license](https://img.shields.io/github/license/iqbaleff214/kamus-banjar-api)](https://github.com/iqbaleff214/kamus-banjar-api/blob/main/LICENSE)
 
</div>

# Kamus Banjar API

Tujuan dari proyek ini adalah membuat API untuk kamus Bahasa Banjar-Indonesia, yang memberikan pengguna kemampuan untuk menerjemahkan kata dari Bahasa Banjar ke Bahasa Indonesia.
Penjelasan secara teori tentang bahasa Banjar silakan merujuk ke halaman [wiki](https://github.com/iqbaleff214/kamus-banjar-api/wiki/Tentang-Bahasa-Banjar).


*Baca dengan bahasa lain: [English](README.en.md).*

## Prasyarat

Proyek ini dibangun menggunakan [**Go version 1.22.2**](https://go.dev/dl/), dan diharapkan untuk dikembangkan menggunakan versi Golang yang serupa untuk mendapatkan hasil sesuai harapan.

## Variabel Lingkungan

| Variabel | Default | Deskripsi |
|---|---|---|
| `PORT` | `:8001` | Port server |
| `MYSQL_DSN` | — | MySQL DSN. **Wajib** untuk mengaktifkan autentikasi, komunitas, kontribusi, dan endpoint admin. Tanpa ini, hanya endpoint kamus publik yang aktif. |
| `JWT_SECRET` | — | Rahasia JWT (min. 32 karakter). Wajib bila `MYSQL_DSN` diset. |
| `JWT_ACCESS_TTL` | `15m` | Masa berlaku access token (format Go duration, contoh: `15m`, `1h`) |
| `JWT_REFRESH_TTL` | `168h` | Masa berlaku refresh token (format Go duration, contoh: `168h`, `7d`) |
| `ADMIN_EMAIL` | — | Email untuk akun admin awal (opsional) |
| `ADMIN_PASSWORD` | — | Password untuk akun admin awal (opsional) |
| `SOURCE_PATH` | `data/` | Path ke file JSON (hanya dengan build tag `fs`) |

> **Catatan:** Endpoint auth, komunitas (write), kontribusi, dan admin **membutuhkan MySQL** (`MYSQL_DSN` harus diset). Tanpa MySQL, hanya endpoint kamus publik (`/alphabets`, `/entries`) yang terdaftar.

## Cara Menjalankan

```shell
# Dengan MySQL
MYSQL_DSN="user:pass@tcp(localhost:3306)/kamusbanjar" \
JWT_SECRET="your-secret-key-at-least-32-chars" \
go run .

# Hanya kamus (tanpa MySQL)
go run .
```

## Cara Menjalankan Tes

Jalankan perintah berikut untuk melakukan test:
```shell
go test -v -cover ./...
```

Untuk integration test (membutuhkan database MySQL test):
```shell
TEST_MYSQL_DSN="user:pass@tcp(localhost:3306)/kamusbanjar_test" \
go test -tags=integration -v ./integration/...
```

## Cara Membuat _File_ Biner

Jalankan perintah berikut untuk membuat _file_ biner:
```shell
go build -ldflags "-s -w" -o ./bin/app .
```

Kemudian jalankan menggunakan perintah `./bin/app`.

## Migrasi Database

Jalankan migrasi secara berurutan:

```shell
mysql -u user -p kamusbanjar < database/migrations/001_schema.sql
mysql -u user -p kamusbanjar < database/seeds/002_seed.sql
mysql -u user -p kamusbanjar < database/migrations/003_refresh_tokens.sql
mysql -u user -p kamusbanjar < database/migrations/004_word_source.sql
mysql -u user -p kamusbanjar < database/migrations/005_contributions.sql
mysql -u user -p kamusbanjar < database/migrations/006_community.sql
```

Atau gunakan Docker Compose (MySQL variant) yang otomatis menjalankan semua migrasi saat pertama kali dijalankan.

## Docker (Quick Start)

**Development dengan MySQL:**
```shell
docker compose -f docker-compose.mysql.yml --env-file .env.mysql up -d
# API tersedia di http://localhost:8001
```

Salin dan isi file env:
```shell
cp .env.mysql.example .env.mysql
# Edit .env.mysql: isi DOMAIN, MYSQL_ROOT_PASSWORD, MYSQL_PASSWORD, JWT_SECRET, dll.
```

**Hanya kamus (embedded JSON, tanpa MySQL):**
```shell
docker compose up -d
```

## Penggunaan

### [GET] /api/v1/alphabets
Mengembalikan semua daftar alfabet.

#### Respon sukses
Respon akan dikembalikan dalam bentuk JSON. Contohnya:
```json
{
  "code": 200,
  "data": [
    {
      "letter": "a",
      "total": 3
    }
  ],
  "message": "All alphabets successfully retrieved.",
  "status": "success"
}
```

#### Respon galat
Respon akan dikembalikan dalam bentuk JSON juga jika terdapat galat. Contohnya:
```json
{
    "code": 500,
    "message": "Internal Server Error",
    "status": "error"
}
```

### [GET] /api/v1/alphabets/{letter}
Mengembalikan daftar kosakata bahasa Banjar berdasarkan alfabet yang diberikan.

#### Parameters
| Name | Keberadaan | Tipe | Deskripsi |
| ----:|:--------:|:----:| ----------- |
| `letter` | wajib | param  | Huruf alfabet |

#### Respon sukses
Respon akan dikembalikan dalam bentuk JSON. Contohnya:
```json
{
  "code": 200,
  "data": {
    "letter": "a",
    "total": 3,
    "words": [
      "abadan",
      "abah",
      "abat"
    ]
  },
  "message": "All words with letter 'a' successfully retrieved.",
  "status": "success"
}
```

### [GET] /api/v1/entries/{word}
Mengembalikan definisi dan arti dari kosakata bahasa Banjar yang diberikan.

#### Parameters
| Name | Keberadaan | Tipe | Deskripsi |
| ----:|:--------:|:----:| ----------- |
| `word` | wajib | param  | kosakata bahasa Banjar. |

#### Respon sukses
Respon akan dikembalikan dalam bentuk JSON. Contohnya:
```json
{
  "code": 200,
  "data": {
    "word": "abah",
    "alphabet": "a",
    "meanings": [
      {
        "definitions": [
          {
            "definition": "ayah",
            "partOfSpeech": "n"
          }
        ]
      }
    ]
  },
  "message": "Definition of word 'abah' successfully retrieved.",
  "status": "success"
}
```

## Daftar Pustaka
- [Departemen Pendidikan Nasional, Pusat Bahasa, Balai Bahasa Banjarmasin. 2008. _Kamus Bahasa Banjar Dialek Hulu-Indonesia_. Banjarbaru.](https://repositori.kemendikdasmen.go.id/2855/1/kamus%20bahasa%20banjar%20dialek%20hulu.pdf)
- [Pusat Pembinaan dan Pengembangan Bahasa, Departemen Pendidikan dan Kebudayaan. 1977. _Kamus Banjar-Indonesia_. Jakarta.](https://repositori.kemendikdasmen.go.id/2888/1/Kamus%20Banjar%20-%20Indonesia%20%20%20%20%20-%20%20%20189h.pdf)

## Lisensi

__Kamus Banjar API__ adalah perangkat lunak _open-source_ yang dilisensikan di bawah lisensi [MIT license](https://github.com/iqbaleff214/kamus-banjar-api/blob/main/LICENSE).
