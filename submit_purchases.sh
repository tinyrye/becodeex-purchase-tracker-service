#!/bin/bash

curl -XPOST -v http://localhost:8999/purchases -d '{ "payer": "DANNON", "points": 300, "timestamp": "2022-10-31T10:00:00Z" }'
curl -XPOST -v http://localhost:8999/purchases -d '{ "payer": "UNILEVER", "points": 200, "timestamp": "2022-10-31T11:00:00Z" }'
curl -XPOST -v http://localhost:8999/purchases -d '{ "payer": "DANNON", "points": -200, "timestamp": "2022-10-31T15:00:00Z" }'
curl -XPOST -v http://localhost:8999/purchases -d '{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2022-11-01T14:00:00Z" }'
curl -XPOST -v http://localhost:8999/purchases -d '{ "payer": "DANNON", "points": 1000, "timestamp": "2022-11-02T14:00:00Z" }'
