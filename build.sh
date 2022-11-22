#!/bin/bash

go mod download
go get purchase-tracker-service
go build -v -o run_http_service
