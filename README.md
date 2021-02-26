# Pasuri

A simple web service for determining has a password been leaked. Similar to the [haveibeenpwned API](https://haveibeenpwned.com) but runnable in an environment where you don't have internet access. Implements [k-anonymity](https://en.wikipedia.org/wiki/K-anonymity).

Pasuri accepts a five character prefix of a SHA1 hashed password and returns their hash suffixes in JSON format. It is left to the user to compare the prefix and suffix to their full password hash.

Example response when queried with **hash?prefix=9D4E1**:

```
[
    "00CEA54681AAC972C1DD81152CB6F840B24",
    "00FCD1A8AAD59272746DDB3E8B0F9151A39",
    "011729341632C79E7484B8EA9D1511DF5B2"
]
```

Pasuri is written in Go and uses SQLite. The motivation for choosing SQLite is it's small size and the simplicity of the data structures. SQLite is good at optimizing integer storage so hashes are persisted as integers instead of hex strings.

## Usage
### Check out and install dependencies
```
git clone https://github.com/mprencipe/pasuri.git
go get
```

### Generate a database (pass.db) from some plaintext password files.
```
./fill-db.sh rockyou.txt xato-net-10-million-passwords-1000000.txt
```

### Run with Docker
```
docker build -t mprencipe/pasuri:0.1 .

docker run --rm -p 8080:8080 --name pasuri mprencipe/pasuri:0.1
```

### Run with Docker
```
go build
./pasuri
```
