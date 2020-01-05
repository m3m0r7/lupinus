#!/bin/sh

docker-compose up
source /etc/systemd/env
/bin/go run main.go
