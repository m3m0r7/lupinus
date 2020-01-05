#!/bin/sh

docker-compose up -d
source /etc/systemd/env
/bin/go run main.go
