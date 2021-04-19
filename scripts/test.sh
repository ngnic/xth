#!/bin/bash

migrate -path migrations -database "$DB_URL" up 

go test -v ./repositories/...