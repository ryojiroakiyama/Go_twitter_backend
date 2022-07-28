# Requirements
- Go
- docker / docker-compose
# Getting Started
## Build & Start  
```bash
#terminal
$ cd [this repository]
$ docker-compose up -d
```
## Access  
access http://localhost:8081/ on Web browser  
### Check API
- open one of API tabs and click `Try it out` and Execute something
- The tabs marked with a key requires `username [your username]` to be entered (simple authentication)
## Stop  
```bash
#terminal
$ docker-compose down
```
# Initialize
```bash
#terminal
docker-compose down
rm -rfd .data/mysql
docker-compose up -d
```
# Log
```
#terminal
docker-compose logs
docker-compose logs -f
docker-compose logs web
docker-compose logs -f web
```
