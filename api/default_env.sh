#!/bin/sh

export TZ="UTC"
export POSTGRES_PORT="11163"
export POSTGRES_DB="room_finder"
export POSTGRES_USER="room_finder"
export POSTGRES_PASSWORD="room_finder_dev_password"
export DATABASE_URL="postgres://room_finder:room_finder_dev_password@127.0.0.1:11163/room_finder?sslmode=disable"
