#!/bin/bash -e
sudo apt update
sudo apt install -y mariadb-server
sudo service mariadb start
sudo mariadb -e "CREATE DATABASE climblive;"
sudo mariadb -e "CREATE USER climblive@localhost IDENTIFIED BY 'secretpassword';"
sudo mariadb -e "GRANT ALL PRIVILEGES ON climblive.* TO climblive@localhost;"

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/gzuidhof/tygo@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install gotest.tools/gotestsum@latest

cd web
pnpm i
cd e2e
pnpm exec playwright install --with-deps
