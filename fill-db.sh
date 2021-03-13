#!/bin/sh

# exit on errors
set -e

usage() {
    echo "Usage: $0 [-h hashfile.txt] [-t textfile.txt]" 1>&2
    exit 1
}

if [ $# -lt 1 ]; then
    usage
    exit 1
fi

files=()

# add file type suffix
while getopts ":h:t:" o; do
    case "${o}" in
        h)
            files+="../${OPTARG}:hash "
            ;;
        t)
            files+="../${OPTARG}:text "
            ;;
        *)
            ;;
    esac
done

# Remove existing database and create new
rm -f pass.db
touch pass.db

# Create table
sqlite3 pass.db 'create table hash (prefix integer not null, part1 integer not null, part2 integer not null, part3 integer not null, primary key(prefix, part1, part2, part3)) without rowid'

# Move database to directory, Go sql library doesn't like relative paths
mv pass.db filldb
cd filldb
go run filldb.go ${files[@]}
mv pass.db ..
cd ..
