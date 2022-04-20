#!/bin/bash

set -eux

# ストリーミングによる確認
docker-compose logs -f

# 確認
# docker-compose logs

# webサーバonly
# docker-compose logs web
# docker-compose logs -f web