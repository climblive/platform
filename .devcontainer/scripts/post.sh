sudo apt update
sudo apt install -y mariadb-server
sudo service mariadb start
sudo mariadb -e "CREATE DATABASE climblive;"
sudo mariadb -e "CREATE USER climblive@localhost IDENTIFIED BY 'secretpassword';"
sudo mariadb -e "GRANT ALL PRIVILEGES ON climblive.* TO climblive@localhost;"
sudo mariadb climblive -e "SOURCE backend/database/scoreboard.sql" --default-character-set utf8mb4
sudo mariadb climblive -e "SOURCE backend/database/samples.sql" --default-character-set utf8mb4

if [[ -n "${CODESPACE_NAME}" ]]; then
    API_URL="https://${CODESPACE_NAME}-8090.app.github.dev"
    sed -i "s,\"API_URL\":.*,\"API_URL\": \"${API_URL}\",g" web/packages/lib/src/config.json
fi