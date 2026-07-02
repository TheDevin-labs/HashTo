# HashTo

![Go Version](https://img.shields.io/badge/go-1.21%2B-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20android%20%7C%20macos%20%7C%20windows-lightgrey)
![Dependencies](https://img.shields.io/badge/dependencies-zero-brightgreen)

*Being a small forge-wrought tool of TheDevinLabs, for the turning of any file into a hash, and of any hash into a form the eye may read.*

---

## I. Of Its Purpose

A hash, once struck, is a fine thing for a machine to compare — yet a poor thing for a person to inspect. HashTo remedies this. It strikes a hash from any file or stream of bytes and seals it in a compact binary vessel (a `.hash` file) — algorithm, digest, timestamp, locale, and the size and name of the source all bound together. From that vessel, HashTo can pour the contents out into **JSON** or **YAML**, plain and legible, and pour them back in again, restoring the original binary form bit for bit.

No third-party library is called upon. The tool is built entirely of the Go standard library, and answers to none but itself.

## II. Of Its Commands

| Flag | Short | Purpose |
|---|---|---|
| `--hash` | | Strike a hash from a file or from stdin |
| `--to-json` | | Pour a `.hash` file out into JSON |
| `--to-yaml` | | Pour a `.hash` file out into YAML |
| `--from-json` | | Pour JSON back into a `.hash` file |
| `--from-yaml` | | Pour YAML back into a `.hash` file |
| `--table-hash FILE` | | Lay every particular of a `.hash` file plain upon the table |
| `--input` | `-i` | The file to read |
| `--output` | `-o` | The file to write (`-` for stdout, on conversions) |
| `--algo` | `-a` | `md5`, `sha1`, `sha256` (default), or `sha512` |
| `--help` | `-h` | Show usage |
| `--version` | `-v` | Show version |

## III. Of Its Use

Strike a hash from a file:

```sh
hashto --hash --input photo.png --algo sha256
```

Strike a hash from whatever comes down the pipe:

```sh
cat notes.txt | hashto --hash --output notes.hash
```

Pour a `.hash` file into YAML for reading, or JSON if that is your custom:

```sh
hashto --to-yaml --input notes.hash
hashto --to-json --input notes.hash --output notes.json
```

Pour it back into binary form, unchanged:

```sh
hashto --from-yaml --input notes.yaml --output rebuilt.hash
```

Lay every particular bare upon the table:

```sh
hashto --table-hash notes.hash
```

```
----------------------
File            : notes.hash
Source          : notes.txt
Source Size     : 8214 bytes
Algorithm       : sha256
Digest (hex)    : 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08
Digest (base64) : n4bQgYhMfWWaL+qgxVrQFaO/TxsrCwiC0Y1sFbDwCgg=
Created         : 2026-07-02T09:14:03Z
Locale/Charset  : en_US.UTF-8
----------------------
```

## IV. Of the `.hash` Format

A `.hash` file is a small binary vessel, laid out thus:

```
"HTOH" magic (4 bytes)
version byte
algorithm byte   (1=md5, 2=sha1, 3=sha256, 4=sha512)
digest length    (uint16) + digest bytes
timestamp        (int64, unix seconds)
locale length    (uint16) + locale bytes
source length    (uint16) + source filename bytes
source size      (int64, bytes read to produce the digest)
```

Nothing here is left to guesswork; the JSON and YAML forms carry the same fields under plain names, and either may be turned back into the binary form without loss.

## V. Of Its Building

```sh
make build      # builds ./hashto for the host platform
make install     # installs to $PREFIX/bin (or /usr/local/bin)
make clean       # sweeps away build artifacts and scratch hashes
```

Cross-compilation needs no ceremony beyond the standard `GOOS`/`GOARCH` of Go itself:

```sh
GOOS=windows GOARCH=amd64 go build -o hashto.exe .
```

## Team 

- CoolyDucks (Abdullah)
