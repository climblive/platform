#!/bin/bash -e
sudo apt update
sudo apt install -y mariadb-server
sudo service mariadb start
sudo mariadb -e "CREATE DATABASE climblive;"
sudo mariadb -e "CREATE USER climblive@localhost IDENTIFIED BY 'secretpassword';"
sudo mariadb -e "GRANT ALL PRIVILEGES ON climblive.* TO climblive@localhost;"
sudo mariadb climblive -e "SOURCE backend/database/climblive.sql" --default-character-set utf8mb4
sudo mariadb climblive -e "SOURCE backend/database/samples.sql" --default-character-set utf8mb4

for i in $(seq -f "%04g" 2 200)
do
    sudo mariadb climblive -e "INSERT INTO contender VALUES (NULL, 1, 1, 'ABCD$i', NULL, NULL, NULL, NULL, FALSE, FALSE)"
done

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/gzuidhof/tygo@latest
go install github.com/pressly/goose/v3/cmd/goose@latest

if [[ -n "${CODESPACE_NAME}" ]]; then
    API_URL="https://${CODESPACE_NAME}-8090.app.github.dev"
    sed -i "s,\"API_URL\":.*,\"API_URL\": \"${API_URL}\",g" web/packages/lib/src/config.json

    sudo apt install -y gh
    gh codespace ports visibility 8090:public -c $CODESPACE_NAME
fi

cd web
pnpm i
cd e2e
pnpm exec playwright install --with-deps