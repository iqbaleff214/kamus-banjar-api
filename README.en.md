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

The goal of this project is to create an API for a Banjar-Indonesian dictionary, providing users with the ability to translate words from Banjar to Indonesian.

*Read this in other languages: [Bahasa Indonesia](README.md).*

## Prerequisite

This project is built using [**Go version 1.22.2**](https://go.dev/dl/), and it is expected to be developed using this specific version of Golang to ensure the desired outcome.

## How to Run

- Install project dependencies using the command `go mod download`.
- Run the service using the command `go run .` or `go run main.go`.

## How to Build

Execute the following command to build the binary:
```shell
go build -ldflags "-s -w" -o ./out .
```

Then you can run the service using the command `./out`.

## Reference
[Pusat Pembinaan dan Pengembangan Bahasa, Departemen Pendidikan dan Kebudayaan. 1977. _Kamus Banjar-Indonesia_. Jakarta.](https://repositori.kemdikbud.go.id/2888/1/Kamus%20Banjar%20-%20Indonesia%20%20%20%20%20-%20%20%20189h.pdf)

## License

__Kamus Banjar API__ is open-sourced software licensed under the [MIT license](https://github.com/iqbaleff214/kamus-banjar-api/blob/main/LICENSE).