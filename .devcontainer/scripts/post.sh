sudo apt update
sudo apt install -y mariadb-server
sudo service mariadb start
sudo mariadb -e "CREATE DATABASE climblive;"
sudo mariadb -e "CREATE USER climblive@localhost IDENTIFIED BY 'climblive';"
sudo mariadb -e "GRANT ALL PRIVILEGES ON climblive.* TO climblive@localhost;"
sudo mariadb climblive -e "SOURCE backend/database/scoreboard.sql"
sudo mariadb climblive -e "SOURCE backend/database/samples.sql"