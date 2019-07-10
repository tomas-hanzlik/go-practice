#!/bin/bash
cd "$(dirname "$0")"

export $(grep -v '^#' .env | xargs)
go run ./main.go