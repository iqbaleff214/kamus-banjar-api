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

## Prasyarat

Proyek ini dibangun menggunakan [**Go version 1.22.2**](https://go.dev/dl/), dan diharapkan untuk dikembangkan menggunakan versi Golang yang serupa untuk mendapatkan hasil sesuai harapan.

## Cara Menjalankan

- Instalasikan dependensi proyek menggunakan perintah `go mod download`.
- Jalankan proyek dengan perintah `go run .` atau `go run main.go`.

## Cara Membuat _File_ Biner

Jalankan perintah berikut untuk membuat _file_ biner:
```shell
go build -ldflags "-s -w" -o ./out .
```

Kemudian jalankan menggunakan perintah `./out`.

## Lisensi

__Kamus Banjar API__ adalah perangkat lunak _open-source_ yang dilisensikan di bawah lisensi [MIT license](https://github.com/iqbaleff214/kamus-banjar-api/blob/main/LICENSE).