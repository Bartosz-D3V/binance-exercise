#!/bin/bash

if ! command -v docker &>/dev/null; then
  echo "Docker not found in the system. Please generate SQLLite file using SQLLite CLI or install Docker."
  exit
fi

echo "Generating SQLLite db file"
docker run --user "$(id -u)":"$(id -g)" \
  -v /"$(pwd)"/data:/db \
  keinos/sqlite3 \
  sh -c "sqlite3 /db/database.db < /db/schema.sql"
echo "Done"
