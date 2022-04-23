#!/bin/bash

echo "tips"
echo "mysql -u yatter -pyatter"
echo "mysql> use yatter;"
echo "mysql> show tables;"
echo "mysql> desc account;"
echo "mysql> show databases;"
echo "mysql> select * from account;"

set -eux
docker-compose exec mysql /bin/bash
