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

## Cara Menjalankan

- Instalasikan dependensi proyek menggunakan perintah `go mod download`.
- Jalankan proyek dengan perintah `go run .` atau `go run main.go`.

## Cara Menjalankan Tes

Jalankan perintah berikut untuk melakukan test:
```shell
go test -v -cover ./...
```

Dengan perintah tersebut, dapat dilihat juga berapa tingkat _coverage_ dari _unit test_ yang ada. Jika melakukan perubahan pada kode, pastikan menjalankan tes terlebih dahulu untuk memastikan semua fungsionalitas tetap bekerja.

## Cara Membuat _File_ Biner

Jalankan perintah berikut untuk membuat _file_ biner:
```shell
go build -ldflags "-s -w" -o ./bin/app .
```

Kemudian jalankan menggunakan perintah `./bin/app`.

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

#### Respon galat
Respon akan dikembalikan dalam bentuk JSON juga jika terdapat galat. Contohnya:
```json
{
  "code": 400,
  "message": "alphabet only has one character",
  "status": "error"
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
          },
          {
            "definition": "mertua laki-laki",
            "partOfSpeech": "n"
          }
        ]
      }
    ],
    "derivatives": [
      {
        "word": "baabah",
        "syllables": "ba.a.bah",
        "definitions": [
          {
            "definition": "berayah; menyebut ayah",
            "partOfSpeech": "v",
            "examples": [
              {
                "bjn": "inya kada baabah",
                "id": "dia tidak berayah"
              }
            ]
          }
        ]
      }
    ]
  },
  "message": "Definition of word 'abah' successfully retrieved.",
  "status": "success"
}
```

#### Respon galat
Respon akan dikembalikan dalam bentuk JSON juga jika terdapat galat. Contohnya:
```json
{
  "code": 404,
  "message": "the word is not found",
  "status": "error"
}
```

## Daftar Pustaka
- [Departemen Pendidikan Nasional, Pusat Bahasa, Balai Bahasa Banjarmasin. 2008. _Kamus Bahasa Banjar Dialek Hulu-Indonesia_. Banjarbaru.](https://repositori.kemdikbud.go.id/2855/1/kamus%20bahasa%20banjar%20dialek%20hulu.pdf)
- [Pusat Pembinaan dan Pengembangan Bahasa, Departemen Pendidikan dan Kebudayaan. 1977. _Kamus Banjar-Indonesia_. Jakarta.](https://repositori.kemdikbud.go.id/2888/1/Kamus%20Banjar%20-%20Indonesia%20%20%20%20%20-%20%20%20189h.pdf)

## Lisensi

__Kamus Banjar API__ adalah perangkat lunak _open-source_ yang dilisensikan di bawah lisensi [MIT license](https://github.com/iqbaleff214/kamus-banjar-api/blob/main/LICENSE).