#!/bin/bash

app_port=$("./scripts/free_app_port.sh")


source $APP_FOLDER/$1/.env

dapr run \
        --app-id "$1" \
        --app-port "$app_port" \
        --app-protocol grpc \
        --components-path "$DAPR_COMPONENTS" \
        --config "$DAPR_CONFIG" \
        --log-level debug \
        go run "./$APP_FOLDER/$1"