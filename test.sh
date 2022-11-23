#!/bin/bash

SCRIPT_DIR=$(dirname $0)
pushd $SCRIPT_DIR > /dev/null

go test

./run_http_service > http_service.log 2> http_service.log &
http_service_pid=$!
sleep 3
# Test Accounts are Created by the Service Startup so purchases
# will line up to the.Accounts
./submit_purchases.sh
./spend_points.sh

is_running=$(ps aux | grep $http_service_pid)
if [[ "$is_running" != "" ]]; then
	kill $http_service_pid
fi