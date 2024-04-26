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


## About the Banjar Language

One of the provinces on the island of Kalimantan is South Kalimantan (Kalsel). Almost the entire area of South Kalimantan is inhabited by Banjar people. The Banjar language (BB) serves as the lingua franca for the Banjar community and functions as their everyday communication tool.

The Banjar language is not only known in the South Kalimantan region but also in the coastal areas of Central Kalimantan (Kalteng) and East Kalimantan (Kaltim), and even in small parts of Sumatra, such as Muara Tungkal, Sapat, and Tambilahan.

The Banjar language has two dialects, namely Banjar Hulu Dialect and Banjar Kuala Dialect. Some phonemes and vocabulary in the Banjar Hulu Dialect (BBDH) bear similarities and resemblances to Indonesian (BI) despite their differing positions.

The similarity in phonemes and vocabulary between BBDH and BI, for example,

| BI | BBDH | English |
| ---- | -- | ------- |
| `lambat` | `lambat` | `slow` |
| `kayu` | `kayu` | `wood` |
| `malam` | `malam` | `night` |
| `makan` | `makan` | `eat` |

The resemblance in phonemes and vocabulary between BBDH and BI, for example,

| BI | BBDH | English |
| ---- | -- | ------- |
| `beri` | `bari` | `give` |
| `hari` | `ari` | `day` |
| `lubang` | `luwang` | `hole` |
| `meja` | `mija` | `table` |

Based on the observations above, both BBDH and BI recognize the vowels [a], [i], and [u]. Additionally, the similarity in expression leads to specific phonemes found in the vocabulary of BBDH and BI, for example, the meaning of `hari` (`day`) is written and pronounced as `hari` in BI while in BBDH it is pronounced as `ari`.

Apart from the points mentioned, there are differences between BBDH and BI as two distinct languages. For instance,

| BI | BBDH | English |
| ---- | -- | ------- |
| `cantik` | `bungas` | `beautiful` |
| `mampu` | `kawa` | `capable` |
| `arah` | `ampah` | `direction` |
| `luas` | `ligar` | `wide` |

Most of the vocabulary found in BI is indeed not found in BBDH, and vice versa. Of course, the differences between BBDH and BI are countless, besides the similarities and resemblances that also cannot be ignored.

### Alphabet and Spelling

Entries in this dictionary are arranged alphabetically, starting from base words, words with affixes, repeated words, and phrases (word combinations). Syllable segmentation is based on the General Guidelines for Indonesian Spelling, Enhanced Edition. It includes some examples of word usage, as well as a brief explanation of the phonology and morphology of BBDH. Letters included in parentheses `(...)` indicate that the letter is not used in BBDH writing. The alphabet is as follows,

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

### Phonemes of the Banjar Language

<table border="1">
  <thead>
    <tr>
      <th rowspan="2"> Number <br /> Types </th>
      <th rowspan="2"> Phonetic <br /> Symbols</th>
      <th rowspan="2"> Spelling </th>
      <th colspan="3">Examples of Use in Three Positions</th>
    </tr>
    <tr>
      <th>Beginning</th>
      <th>Middle</th>
      <th>End</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td colspan="2">Vowel</td>
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
      <td colspan="2">Diphthong Vowels</td>
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
      <td colspan="2">Consonant</td>
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

### Possessive Form

The term 'possessive form' here refers to the syllable structure in BBDH. The possessive form of base entries, affixes, or reduplication is as follows:

- Base Entry
  
  `V` = Vokal (Vowel) | `K` = Konsonan (Consonant)

  Having two syllables:
  
  | Arrangement | Word |
  | ------- | ---- |
  | VKV | una 'sentuh' |
  | VKVK | amun 'jika' |
  | KVKV | saku 'mungkin' |
  | KVKVK | tagal 'tetapi' |
  | KVKKVK | limbui 'kuyup' |

  Having three syllables:

  | Arrangement | Word |
  | ------- | ---- |
  | KVKVKV | wahini 'sekarang' |
  | KVKVKVK | karukut 'cakar' |
  | KVKVKKV | kalambu 'kelambu' |
  | KVKKVKV | hintalu 'telur' |
  | KVKKVKVK | tantaran 'jorang' |
  | KVKVKKVK | karamput 'bohong' |
  | KVKKVKVK | bungkalang 'keranjang buah' |

- Affixes

  The affix forms in BBDH consist of `KV-`, and `-KVK`, `-VK`, `-V`, `-VKVK` such as the prefixes `ma-`, `ba-`, `ta-`, `sa-`, `ka-`, `pa-`, for example:

  | Affix | Base Word | Affixed Word |
  | --- | --- | --- |
  | `ma-` | bari 'beri' | mambari 'memberi' |
  | `ba-` | jariji 'jemari' | bajariji 'berjemari' |
  | `ta-` | rumpak 'tabrak' | tarumpak 'tertabrak' |
  | `sa-` | bujur 'benar' | sabujurnya 'sebenarnya' |
  | `ka-` | handak 'pendek' | kehandapan 'kependekan' |
  | `pa-` | uncit 'terakhir' | pauncitnya 'paling terakhir' |

  The prefix `ka-an` signifies superlative meaning, for example:

  | Affix | Base Word | Affixed Word |
  | --- | --- | --- |
  | `ka-an` | handak 'ingin' | kahandakan 'terlalu ingin' |

  `ta-` can attach to verbs, adjectives, and nouns, for example

  | Affix | Base Word | Affixed Word |
  | --- | --- | --- |
  | `ta-` | pukul 'pukul' | tapukul 'terpukul' |
  | `ta-` | bungas 'cantik' | tabungas 'lebih cantik' |
  | `ta-` | unjun 'kail' | taunjun 'terpancing' |
  | `ta-` | akal 'akal' | taakal 'lebih berakal' |

  `pa-`+`-nya` and `pa-`+`nasal`+`-nya` both signify a superlative meaning, for example

  | Affix | Base Word | Affixed Word |
  | --- | --- | --- |
  | `pa-`+`-nya` | lambat 'lambat' | palambatnya 'paling lambat' |
  | `pa-`+`-nya` | pintar 'pintar' | pamintarnya 'paling pintar' |
  | `pa-`+`-nya` | dikit 'dikit' | pandikitnya 'paling sedikit' |
  | `pa-`+`-nya` | sabak 'berantakan' | panyabaknya 'paling berantakan' |
  | `pa-`+`-nya` | kurus 'kurus' | pangurusnya 'paling kurus' |

  and such as the suffixes: `-an`, `-i`, `-akan`, for example

  | Affix | Base Word | Affixed Word |
  | --- | --- | --- |
  | `-an` | bulik 'pulang' | bulikan 'pada pulang' |
  | `-i` | parak 'dekat' | paraki 'dekati' |
  | `-akan` | surui 'sisir' | suruiakan 'sisirkan' |

- Reduplication

  The reduplication form in BBDH takes the form of _dwipurwa_ or the repetition of part or all of the initial syllable of a word, for example

  | Base Word | Reduplication |
  | ---------- | ----------- |
  | ranai 'diam' | raranai 'diam-diam' |
  | handap 'pendek' | hahandap 'pendek-pendek' |

  In addition, there are some that must end with -an, thus having a plural meaning, for example

  | Base Word | Reduplication |
  | ---------- | ----------- |
  | umpat 'ikut' | uumpatan 'ikut-ikutan' |
  | ingat 'ingat' | iingatan 'yang diingat-ingat' |

### Abbreviation

| Abbreviation | Expansion    |
| ------------ | ------------ |
| a            | adjektiva    |
| n            | nomina       |
| pro          | pronomina    |
| adv          | adverbia     |
| v            | verba        |

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

## Usage

### [GET] /api/v1/alphabets
Returning a list of alphabet information.

#### Success response
The response will be returned an JSON. For example:
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
The response will be returned as JSON in case of an error. For example:
```json
{
    "code": 500,
    "message": "Internal Server Error",
    "status": "error"
}
```

### [GET] /api/v1/alphabets/{letter}
Returning a list of Banjar words according to the given letter.

#### Parameters
| Name | Required | Type | Description |
| ----:|:--------:|:----:| ----------- |
| `letter` | required | param  | Alphabetic letter |

#### Success response
The response will be returned an JSON. For example:
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

#### Success Response Field
| JSON Field | Type | Description |
| :---|:----:| ----------- |
| `data.meanings[i].definitions[i].refer` | string  | Related words |
| `data.meanings[i].definitions[i].levelOfPoliteness` | string  | The level of politeness of a word. Returns a value of `0` if the word is __neutral__, `1` if the word is __polite__, and will return `-1` if the word is categorized as __rude__. |

#### Error response
The response will be returned as JSON in case of an error. For example:
```json
{
  "code": 400,
  "message": "alphabet only has one character",
  "status": "error"
}
```

### [GET] /api/v1/entries/{word}
Returning the definition and meaning according to the given Banjar word.

#### Parameters
| Name | Required | Type | Description |
| ----:|:--------:|:----:| ----------- |
| `word` | required | param  | Banjar word. |

#### Success response
The response will be returned an JSON. For example:
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

#### Success Response Field
| JSON Field | Type | Description |
| :---|:----:| ----------- |
| `data.meanings[i].definitions[i].refer` | string  | Related words |
| `data.meanings[i].definitions[i].levelOfPoliteness` | string  | The level of politeness of a word. Returns a value of `0` if the word is __neutral__, `1` if the word is __polite__, and will return `-1` if the word is categorized as __rude__. |

#### Error response
The response will be returned as JSON in case of an error. For example:
```json
{
  "code": 404,
  "message": "the word is not found",
  "status": "error"
}
```

## Reference
- [Departemen Pendidikan Nasional, Pusat Bahasa, Balai Bahasa Banjarmasin. 2008. _Kamus Bahasa Banjar Dialek Hulu-Indonesia_. Banjarbaru.](https://repositori.kemdikbud.go.id/2855/1/kamus%20bahasa%20banjar%20dialek%20hulu.pdf)
- [Pusat Pembinaan dan Pengembangan Bahasa, Departemen Pendidikan dan Kebudayaan. 1977. _Kamus Banjar-Indonesia_. Jakarta.](https://repositori.kemdikbud.go.id/2888/1/Kamus%20Banjar%20-%20Indonesia%20%20%20%20%20-%20%20%20189h.pdf)

## License

__Kamus Banjar API__ is open-sourced software licensed under the [MIT license](https://github.com/iqbaleff214/kamus-banjar-api/blob/main/LICENSE).