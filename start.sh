#!/bin/sh

# Change working directory
cd /var/www/lupinus

while :
  do
    echo "Waiting for upon docker."
    docker ps 1>/dev/null 2>&1

    if [ $? == 0 ]; then
      echo "Docker is running"
      break
    fi
    echo "."
    sleep 1
  done

docker-compose up -d
source /etc/systemd/env
/bin/go run main.go
