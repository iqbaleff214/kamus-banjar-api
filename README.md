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


*Baca dengan bahasa lain: [English](README.en.md).*


## Tentang Bahasa Banjar

Salah satu provinsi di pulau Kalimantan adalah Kalimantan Selatan (Kalsel). Hampir seluruh wilayah Kalsel dihuni oleh orang Banjar. Bahasa Banjar (BB) bagi masyarakat Banjar merupakan bahasa pengantar yang berfungsi sebagai alat komunikasi sehari-hari.

Bahasa Banjar dalam penyebarannya tidak hanya dikenal di wilayah Kalsel saja, tetapi juga di pesisir Kalimantan Tengah (Kalteng) dan Kalimantan Timur (Kaltim) bahkan sampai di sebagian kecil daerah Sumatera, seperti Muara Tungkal, Sapat, dan Tambilahan.

Bahasa Banjar memiliki dua dialek, yaitu Banjar Dialek Hulu dan Banjar Kuala. Ada sebagian fonem maupun kosakata Bahasa Banjar Dialek Hulu (BBDH) yang memiliki persamaan dan kemiripan dengan Bahasa Indonesia (BI) meski kedudukannya berbeda.

Persamaan fonem dan kosakata antara BBDH dan BI, contohnya,

| BI | BBDH |
| ---- | -- |
| `lambat` | `lambat` |
| `kayu` | `kayu` |
| `malam` | `malam` |
| `makan` | `makan` |

Kemiripan fonem dan kosakata antara BBDH dan BI, contohnya,

| BI | BBDH |
| ---- | -- |
| `beri` | `bari` |
| `hari` | `ari` |
| `lubang` | `luwang` |
| `meja` | `mija` |

Berdasarkan pengamatan di atas, antara BBDH dan BI sama-sama mengenal vokal [a], [i], dan [u]. Selain itu, kemiripan dari segi pengungkapan mengarah kepada fonem tertentu yang terdapat pada kosakata BBDH dan BI, contohnya makna `hari` yang dalam BI tulisan dan pengungkapannya `hari` sedang dalam BBDH `ari`.

Selain hal-hal yang telah dikemukakan, antara BBDH dan BI sebagai dua buah bahasa memiliki perbedaan, contohnya,

| BI | BBDH |
| ---- | -- |
| `cantik` | `bungas` |
| `mampu` | `kawa` |
| `arah` | `ampah` |
| `luas` | `ligar` |

Sebagian besar kosakata yang terdapat dalam BI memang tidak terdapat dalam BBDH begitu pula sebaliknya. Tentu saja perbedaan antara BBDH dan BI ini tidak terhitung banyaknya selain persamaan dan kemiripan yang juga tidak bisa diindahkan keberadaannya.

### Abjad dan Ejaan

Penulisan entri dalam kamus ini disusun secara alfabetis, berurutan dari kata dasar, kata berimbuhan, kata berulang, dan frasa (gabungan kata). Pemenggalan suku kata berdasarkan Pedoman Umum Ejaan Bahasa Indonesia yang Disempurnakan. Ditambah dengan sebagian contoh penggunaan kosakatanya, serta penjelasan secara singkat mengenai fonologi dan morfologi BBDH. Bagi huruf yang dimasukkan ke dalam tanda kurung `(...)` menunjukkan bahwa huruf tersebut tidak dipakai dalam penulisan BBDH. Adapun abjadnya sebagai berikut,

- `a`
- `b`
- `c`
- `d`
- `(e)`
- `(f)`
- `g`
- `h`
- `i`
- `j`
- `k`
- `l`
- `m`
- `n/ng-ny`
- `(o)`
- `p`
- `(q)`
- `r`
- `s`
- `t`
- `u`
- `(v)`
- `w`
- `y`
- `(z)`

### Fonem Bahasa Banjar

<table>
  <thead>
    <tr>
      <th rowspan="2"> Jenis <br /> Nomor </th>
      <th rowspan="2"> Simbol <br /> Fonetis</th>
      <th rowspan="2"> Ejaan </th>
      <th colspan="3">Contoh Pemakaian Dalam Tiga Posisi</th>
    </tr>
    <tr>
      <th>Awal</th>
      <th>Tengah</th>
      <th>Akhir</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td colspan="2">Vokal</td>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>1</td>
      <td>[a]</td>
      <td>a</td>
      <td><b>a</b>but  `ribut`</td>
      <td>ba<b>a</b>ng `azan`</td>
      <td>tatamb<b>a</b> `obat`</td>
    </tr>
    <tr>
      <td>2</td>
      <td>[i]</td>
      <td>i</td>
      <td><b>i</b>suk  `besok`</td>
      <td>k<b>i</b>pit `sempit`</td>
      <td>jarij<b>i</b> `jemari`</td>
    </tr>
    <tr>
      <td>3</td>
      <td>[u]</td>
      <td>u</td>
      <td><b>u</b>mpat  `ikut`</td>
      <td>j<b>u</b>hut `tarik`</td>
      <td>sang<b>u</b> `bekal`</td>
    </tr>
    <tr>
      <td colspan="2">Vokal Rangkap</td>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>4</td>
      <td>[au]</td>
      <td>aw</td>
      <td>-</td>
      <td>-</td>
      <td>sil<b>au</b> `silau`</td>
    </tr>
    <tr>
      <td>5</td>
      <td>[ai]</td>
      <td>ay</td>
      <td>-</td>
      <td>s<b>ai</b>tan `setan`</td>
      <td>wad<b>ai</b> `kue`</td>
    </tr>
    <tr>
      <td>6</td>
      <td>[ui]</td>
      <td>uy</td>
      <td>-</td>
      <td>-</td>
      <td>rup<b>ui</b> `keropos`</td>
    </tr>
    <tr>
      <td colspan="2">Konsonan</td>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>7</td>
      <td>[p]</td>
      <td>p</td>
      <td><b>p</b>ayu `laku`</td>
      <td>su<b>p</b>an `malu`</td>
      <td>suma<b>p</b> `kukus`</td>
    </tr>
    <tr>
      <td>8</td>
      <td>[b]</td>
      <td>b</td>
      <td><b>b</b>alu `janda`</td>
      <td>hu<b>p</b>an `uban`</td>
      <td>-</td>
    </tr>
    <tr>
      <td>9</td>
      <td>[t]</td>
      <td>t</td>
      <td><b>t</b>alu `tiga`</td>
      <td>ma<b>t</b>an `dari`</td>
      <td>kipi<b>t</b> `sempit`</td>
    </tr>
    <tr>
      <td>10</td>
      <td>[d]</td>
      <td>d</td>
      <td><b>d</b>awa `tuduh`</td>
      <td>pa<b>d</b>u `dapur`</td>
      <td>-</td>
    </tr>
    <tr>
      <td>11</td>
      <td>[c]</td>
      <td>c</td>
      <td><b>c</b>aluk `rogoh`</td>
      <td>pa<b>c</b>ang `bakal`</td>
      <td>-</td>
    </tr>
    <tr>
      <td>12</td>
      <td>[j]</td>
      <td>j</td>
      <td><b>j</b>angkau `raih`</td>
      <td>ja<b>j</b>ak `injak</td>
      <td>-</td>
    </tr>
    <tr>
      <td>13</td>
      <td>[k]</td>
      <td>k</td>
      <td><b>k</b>atup `tutup`</td>
      <td>la<b>k</b>as `cepat`</td>
      <td>pata<b>k</b> `kubur`</td>
    </tr>
    <tr>
      <td>14</td>
      <td>[g]</td>
      <td>g</td>
      <td><b>g</b>ampir `kembar`</td>
      <td>a<b>g</b>a `gagap`</td>
      <td>-</td>
    </tr>
    <tr>
      <td>15</td>
      <td>[m]</td>
      <td>m</td>
      <td><b>m</b>uar `benci`</td>
      <td>a<b>m</b>pun `milik`</td>
      <td>anu<b>m</b> `muda`</td>
    </tr>
    <tr>
      <td>16</td>
      <td>[n]</td>
      <td>n</td>
      <td><b>n</b>ahap `mantap`</td>
      <td>ka<b>n</b>as `nanas`</td>
      <td>karanga<b>n</b> `pasir`</td>
    </tr>
    <tr>
      <td>17</td>
      <td>[ŋ]</td>
      <td>ng</td>
      <td><b>ng</b>alih `susah`</td>
      <td>ta<b>ng</b>guh `tebak`</td>
      <td>ladi<b>ng</b> `pisau`</td>
    </tr>
    <tr>
      <td>18</td>
      <td>[ɲ]</td>
      <td>ny</td>
      <td><b>ny</b>iru `nyiru`</td>
      <td>ha<b>ny</b>ar `baru`</td>
      <td>-</td>
    </tr>
    <tr>
      <td>19</td>
      <td>[s]</td>
      <td>s</td>
      <td><b>s</b>urui `sisir`</td>
      <td>ra<b>s</b>uk `cocok`</td>
      <td>bati<b>s</b> `kaki`</td>
    </tr>
    <tr>
      <td>20</td>
      <td>[h]</td>
      <td>h</td>
      <td><b>h</b>alar `sayap`</td>
      <td>mu<b>h</b>a `muka`</td>
      <td>lapa<b>h</b> `capek`</td>
    </tr>
    <tr>
      <td>21</td>
      <td>[l]</td>
      <td>l</td>
      <td><b>l</b>icak `becek`</td>
      <td>ta<b>l</b>ah `habis`</td>
      <td>gana<b>l</b> `besar`</td>
    </tr>
    <tr>
      <td>22</td>
      <td>[r]</td>
      <td>r</td>
      <td><b>r</b>umbis `bocor`</td>
      <td>sa<b>r</b>ak `cerai`</td>
      <td>kita<b>r</b> `geser`</td>
    </tr>
    <tr>
      <td>23</td>
      <td>[w]</td>
      <td>w</td>
      <td><b>w</b>adah `tempat`</td>
      <td>ka<b>w</b>a `mampu`</td>
      <td>-</td>
    </tr>
    <tr>
      <td>24</td>
      <td>[y]</td>
      <td>y</td>
      <td><b>y</b>uta `juta`</td>
      <td>u<b>y</b>ah `garam`</td>
      <td>-</td>
    </tr>
  </tbody>
</table>

### Bentuk Persukuan

Maksud bentuk persukuan di sini adalah struktur suku kata (silabel) dalam BBDH. Bentuk persukuan entri pokok, imbuhan atau reduplikasi adalah sebagai berikut:

- Entri Pokok
  
  `V` = Vokal | `K` = Konsonan

  Yang bersuku dua:
  
  | Susunan | Kata |
  | ------- | ---- |
  | VKV | una 'sentuh' |
  | VKVK | amun 'jika' |
  | KVKV | saku 'mungkin' |
  | KVKVK | tagal 'tetapi' |
  | KVKKVK | limbui 'kuyup' |

  Yang bersuku tiga:

  | Susunan | Kata |
  | ------- | ---- |
  | KVKVKV | wahini 'sekarang' |
  | KVKVKVK | karukut 'cakar' |
  | KVKVKKV | kalambu 'kelambu' |
  | KVKKVKV | hintalu 'telur' |
  | KVKKVKVK | tantaran 'jorang' |
  | KVKVKKVK | karamput 'bohong' |
  | KVKKVKVK | bungkalang 'keranjang buah' |

- Imbuhan

  Bentuk imbuhan BBDH terdiri atas `KV-`, dan `-KVK`, `-VK`, `-V`, `-VKVK` seperti awalan `ma-`, `ba-`, `ta-`, `sa-`, `ka-`, `pa-`, misalnya:

  | Imbuhan | Kata Dasar | Kata Berimbuhan |
  | --- | --- | --- |
  | `ma-` | bari 'beri' | mambari 'memberi' |
  | `ba-` | jariji 'jemari' | bajariji 'berjemari' |
  | `ta-` | rumpak 'tabrak' | tarumpak 'tertabrak' |
  | `sa-` | bujur 'benar' | sabujurnya 'sebenarnya' |
  | `ka-` | handak 'pendek' | kehandapan 'kependekan' |
  | `pa-` | uncit 'terakhir' | pauncitnya 'paling terakhir' |

  Imbuhan `ka-an` memiliki makna superlatif, misalnya:

  | Imbuhan | Kata Dasar | Kata Berimbuhan |
  | --- | --- | --- |
  | `ka-an` | handak 'ingin' | kahandakan 'terlalu ingin' |

  `ta-` bisa melekati verba, adjektifa, dan nomina, misalnya

  | Imbuhan | Kata Dasar | Kata Berimbuhan |
  | --- | --- | --- |
  | `ta-` | pukul 'pukul' | tapukul 'terpukul' |
  | `ta-` | bungas 'cantik' | tabungas 'lebih cantik' |
  | `ta-` | unjun 'kail' | taunjun 'terpancing' |
  | `ta-` | akal 'akal' | taakal 'lebih berakal' |

  `pa-`+`-nya` dan `pa-`+`nasal`+`-nya` memiliki makna superlatif, misalnya

  | Imbuhan | Kata Dasar | Kata Berimbuhan |
  | --- | --- | --- |
  | `pa-`+`-nya` | lambat 'lambat' | palambatnya 'paling lambat' |
  | `pa-`+`-nya` | pintar 'pintar' | pamintarnya 'paling pintar' |
  | `pa-`+`-nya` | dikit 'dikit' | pandikitnya 'paling sedikit' |
  | `pa-`+`-nya` | sabak 'berantakan' | panyabaknya 'paling berantakan' |
  | `pa-`+`-nya` | kurus 'kurus' | pangurusnya 'paling kurus' |

  dan seperti akhiran: `-an`, `-i`, `-akan`, misalnya

  | Imbuhan | Kata Dasar | Kata Berimbuhan |
  | --- | --- | --- |
  | `-an` | bulik 'pulang' | bulikan 'pada pulang' |
  | `-i` | parak 'dekat' | paraki 'dekati' |
  | `-akan` | surui 'sisir' | suruiakan 'sisirkan' |

- Reduplikasi

  Bentuk reduplikasi dalam BBDH berbentuk dwipurwa atau pengulangan sebagian atau seluruh suku kata awal sebuah kata, misalnya

  | Kata Dasar | Reduplikasi |
  | ---------- | ----------- |
  | ranai 'diam' | raranai 'diam-diam' |
  | handap 'pendek' | hahandap 'pendek-pendek' |

  selain itu ada sebagian yang harus berakhiran `-an`, sehingga memiliki makna jamak, misalnya  

  | Kata Dasar | Reduplikasi |
  | ---------- | ----------- |
  | umpat 'ikut' | uumpatan 'ikut-ikutan' |
  | ingat 'ingat' | iingatan 'yang diingat-ingat' |


### Kategori Kelas Kata

| **Singkatan** | **Kepanjangan**          | **Penjelasan**                                                   |
|---------------|--------------------------|-------------------------------------------------------------------|
| **n**         | Nomina                  | Kata benda; merujuk pada orang, tempat, atau benda.              |
| **v**         | Verba                   | Kata kerja; menyatakan tindakan, proses, atau keadaan.           |
| **a**         | Adjektiva               | Kata sifat; menerangkan atau membatasi nomina.                   |
| **pro**       | Pronomina               | Kata ganti; menggantikan nomina.                                 |
| **adv**       | Adverbia                | Kata keterangan; menerangkan verba, adjektiva, atau kalimat.     |
| **num**       | Numeralia               | Kata bilangan; menunjukkan jumlah atau urutan.                   |
| **p**         | Partikel                | Kata tugas; tidak berubah bentuk (misalnya: "lah", "kah").       |
| **konj**      | Konjungsi               | Kata sambung; menghubungkan klausa, frasa, atau kata.            |
| **prep**      | Preposisi               | Kata depan; menunjukkan hubungan antara kata (contoh: "di", "ke").|
| **interj**    | Interjeksi              | Kata seru; mengekspresikan emosi atau reaksi spontan.            |
| **klit**      | Klitika                 | Kata pendek yang melekat pada kata lain (contoh: "-ku", "-mu").  |
| **dem**       | Demonstrativa           | Kata tunjuk; menunjuk sesuatu (contoh: "ini", "itu").            |
| **art**       | Artikel                 | Kata sandang; membatasi nomina (contoh: "sang", "si").           |



## Prasyarat

Proyek ini dibangun menggunakan [**Go version 1.22.2**](https://go.dev/dl/), dan diharapkan untuk dikembangkan menggunakan versi Golang yang serupa untuk mendapatkan hasil sesuai harapan.

## Cara Menjalankan

- Instalasikan dependensi proyek menggunakan perintah `go mod download`.
- Jalankan proyek dengan perintah `go run .` atau `go run main.go`.

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

#### Error response
Respon akan dikembalikan dalam bentuk JSON juga jika terdapat eror. Contohnya:
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

#### Error response
Respon akan dikembalikan dalam bentuk JSON juga jika terdapat eror. Contohnya:
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

#### Error response
Respon akan dikembalikan dalam bentuk JSON juga jika terdapat eror. Contohnya:
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