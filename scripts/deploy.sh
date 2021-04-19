#!/bin/bash

migrate -path migrations -database "$DB_URL" up 

go build -o github-app

exec /opt/app/github-app