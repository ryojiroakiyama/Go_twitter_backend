#!/bin/bash

set -eux

docker-compose exec mysql /bin/bash
## mysqlに入る
#mysql -u yatter -pyatter

## yatterデータベースを使用
#mysql> use yatter;

## テーブル表示
#mysql> show tables;
## accountテーブルのカラム一覧を表示
#mysql> desc account;
## データベースの確認
#mysql> show databases;
## accountテーブルのレコードを表示
#mysql> select * from account;