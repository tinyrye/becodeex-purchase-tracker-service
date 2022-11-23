#!/bin/bash

curl -XPOST -H 'Content-Type: application/json' http://localhost:8999/rewards/spend -d '{"points": 5000}' | jq .
