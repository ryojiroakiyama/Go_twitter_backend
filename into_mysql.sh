#!/bin/bash

echo "tips"
echo "mysql> use yatter;"
echo "mysql> show tables;"
echo "mysql> desc account;"
echo "mysql> show databases;"
echo "mysql> select * from account;"

set -eux
echo "mysql -u yatter -pyatter" | pbcopy # コンテナに入った後, コピペでmysqlまで入れるように
docker-compose exec mysql /bin/bash
