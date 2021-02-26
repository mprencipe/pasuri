#!/bin/sh

if [ $# -lt 1 ]; then
    echo "At least one password file expected."
    exit 1
fi

# Prefix files with relative path
files=("$@") 
for i in "${!files[@]}"
do
    files[i]="../${files[i]}"
done

# Remove existing database and create new
rm -f pass.db
touch pass.db

# Create table
sqlite3 pass.db 'create table hash (prefix text not null, suffix text not null, primary key(prefix, suffix))'

# Move database to directory, Go sql library doesn't like relative paths
mv pass.db filldb
cd filldb
go run filldb.go "${files[@]}"
mv pass.db ..
cd ..
