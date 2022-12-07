#!/bin/sh
docker run -d --rm -p 8000:8000 surrealdb/surrealdb:latest start --log warn --user root --pass root
go run .