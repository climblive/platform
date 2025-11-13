#!/bin/bash

for i in $(seq -f "%04g" 1 200)
do
    sudo mariadb climblive -e "INSERT INTO contender VALUES (NULL, 1, 1, 'ABCD$i', NULL, NULL, NULL, NULL, FALSE, FALSE)"
done