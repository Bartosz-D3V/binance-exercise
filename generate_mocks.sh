#!/bin/bash

if ! command -v mockgen &>/dev/null; then
  echo "mockgen not found. Please install it using Go mod: https://github.com/golang/mock"
  exit
fi

echo "Generating mocks using mockgen"
mockgen -source ./app/transaction/repository.go -destination ./app/mock/repository.go -package mock
echo "Done"
