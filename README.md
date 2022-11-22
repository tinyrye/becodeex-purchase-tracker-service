# Purchase Rewards Tracker Service #

## Purpose ##

- To implement an exercise for a Purchase Rewards Application.
- Offer registing and tracking of transactions and rewards for Users.

## Actors ##

1. Purchasers - these are the users that purchase items from Payers/Partners in order to accumulate rewards from the Payers.
2. Payers/Partners - these are the vendors/businesses that sell products to Purchasers and which pay rewards to the Purchasers as a result of purchases.

## Services ##

### Manage Reward Progress/Balances ###

1. Add/List Payers
2. Obtain Payer Balance per Purchaser.
3. Observe Purchase Transaction of a Purchaser.

## How to Use ##

### Prerequisite ###

A Linux or Unix like system with `golang` installed at 1.18

Optional prerequisite of `jq` [JQ](https://stedolan.github.io/jq) is used in example to format responses.

### Building ###

Execute `build.sh`.  This builds a `run_http_service` executable.

### Sample Execution ###

First, startup the server via `run_http_service`

The logs should show, 

```
(base) [littleking@fedora purchase-tracker-service]$ ./run_http_service 
2022/11/22 14:46:29 Added account DANNON with tokens [d da dan dann]
2022/11/22 14:46:29 Added account UNILEVER with tokens [u un uni unil unile unilev]
2022/11/22 14:46:29 Added account MILLER COORS with tokens [  m mi mil mill mille miller miller  miller c miller co miller coo]
2022/11/22 14:46:29 Listening with HTTP server on :8999
```

Next, list current balances `curl -XGET http://localhost:8999/payers/balances`

```json
[
  {
    "payer": {
      "Id": "UNILEVER",
      "Name": "Unilever",
      "CreationTimestamp": "2022-11-22T14:35:34.183517886-06:00"
    },
    "Points": 0
  },
  {
    "payer": {
      "Id": "MILLER COORS",
      "Name": "Miller Coors",
      "CreationTimestamp": "2022-11-22T14:35:34.18356334-06:00"
    },
    "Points": 0
  },
  {
    "payer": {
      "Id": "DANNON",
      "Name": "Dannon",
      "CreationTimestamp": "2022-11-22T14:35:34.1833076-06:00"
    },
    "Points": 0
  }
]
```

Next, add a sample purchase to DANNON like `curl -XPOST -v http://localhost:8999/purchases -d '{"payer": "DANNON", "points": 101}'`

The API responds with

```json
{
    "payer": {
        "id": "DANNON",
        "name": "Dannon",
        "creationTimestamp": "2022-11-22T14:42:07.499412795-06:00"
    },
    "points": 101
```

This just verifies that the service is working properly.

Now, restart the service And supply the real test load:

```bash
./submit_purchases.sh
````

The output should look a lot like

```
(base) [littleking@fedora purchase-tracker-service]$ ./submit_purchases.sh 
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8999...
* Connected to localhost (127.0.0.1) port 8999 (#0)
> POST /purchases HTTP/1.1
> Host: localhost:8999
> User-Agent: curl/7.82.0
> Accept: */*
> Content-Length: 73
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 22 Nov 2022 20:49:30 GMT
< Content-Length: 113
< Content-Type: text/plain; charset=utf-8
< 
{"payer":{"id":"DANNON","name":"Dannon","creationTimestamp":"2022-11-22T14:46:29.053044665-06:00"},"points":300}
* Connection #0 to host localhost left intact
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8999...
* Connected to localhost (127.0.0.1) port 8999 (#0)
> POST /purchases HTTP/1.1
> Host: localhost:8999
> User-Agent: curl/7.82.0
> Accept: */*
> Content-Length: 75
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 22 Nov 2022 20:49:30 GMT
< Content-Length: 117
< Content-Type: text/plain; charset=utf-8
< 
{"payer":{"id":"UNILEVER","name":"Unilever","creationTimestamp":"2022-11-22T14:46:29.053265166-06:00"},"points":200}
* Connection #0 to host localhost left intact
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8999...
* Connected to localhost (127.0.0.1) port 8999 (#0)
> POST /purchases HTTP/1.1
> Host: localhost:8999
> User-Agent: curl/7.82.0
> Accept: */*
> Content-Length: 74
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 22 Nov 2022 20:49:30 GMT
< Content-Length: 113
< Content-Type: text/plain; charset=utf-8
< 
{"payer":{"id":"DANNON","name":"Dannon","creationTimestamp":"2022-11-22T14:46:29.053044665-06:00"},"points":100}
* Connection #0 to host localhost left intact
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8999...
* Connected to localhost (127.0.0.1) port 8999 (#0)
> POST /purchases HTTP/1.1
> Host: localhost:8999
> User-Agent: curl/7.82.0
> Accept: */*
> Content-Length: 81
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 22 Nov 2022 20:49:30 GMT
< Content-Length: 127
< Content-Type: text/plain; charset=utf-8
< 
{"payer":{"id":"MILLER COORS","name":"Miller Coors","creationTimestamp":"2022-11-22T14:46:29.053319156-06:00"},"points":10000}
* Connection #0 to host localhost left intact
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8999...
* Connected to localhost (127.0.0.1) port 8999 (#0)
> POST /purchases HTTP/1.1
> Host: localhost:8999
> User-Agent: curl/7.82.0
> Accept: */*
> Content-Length: 74
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 22 Nov 2022 20:49:30 GMT
< Content-Length: 114
< Content-Type: text/plain; charset=utf-8
< 
{"payer":{"id":"DANNON","name":"Dannon","creationTimestamp":"2022-11-22T14:46:29.053044665-06:00"},"points":1100}
```

And the points balances from `curl -XGET http://localhost:8999/payers/balances` should be

```bash
(base) [littleking@fedora purchase-tracker-service]$ curl -XGET http://localhost:8999/payers/balances | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   360  100   360    0     0   674k      0 --:--:-- --:--:-- --:--:--  351k
[
  {
    "payer": {
      "id": "DANNON",
      "name": "Dannon",
      "creationTimestamp": "2022-11-22T14:46:29.053044665-06:00"
    },
    "points": 1100
  },
  {
    "payer": {
      "id": "UNILEVER",
      "name": "Unilever",
      "creationTimestamp": "2022-11-22T14:46:29.053265166-06:00"
    },
    "points": 200
  },
  {
    "payer": {
      "id": "MILLER COORS",
      "name": "Miller Coors",
      "creationTimestamp": "2022-11-22T14:46:29.053319156-06:00"
    },
    "points": 10000
  }
]
````

Now, let's spend 5000 Points per the Example:

```bash
curl -XPOST -H 'Content-Type: application/json' http://localhost:8999/rewards/spend -d '{"points": 5000}'
```

The resulting allocation of Points from the Payers is as expected:

```json
[
    {
        "payer": {
            "id": "DANNON",
            "name": "Dannon",
            "creationTimestamp": "2022-11-22T15:21:55.433312265-06:00"
        },
        "points": 1000
    },
    {
        "payer": {
            "id": "UNILEVER",
            "name": "Unilever",
            "creationTimestamp": "2022-11-22T15:21:55.433540934-06:00"
        },
        "points": 0
    },
    {
        "payer": {
            "id": "MILLER COORS",
            "name": "Miller Coors",
            "creationTimestamp": "2022-11-22T15:21:55.433596051-06:00"
        },
        "points": 5300
    }
]
```

Which matches what the exercise deems is the correct response.

The server's log from the entire process is

```
(base) [littleking@fedora purchase-tracker-service]$ ./run_http_service 
2022/11/22 15:21:55 Added account DANNON with tokens [d da dan dann]
2022/11/22 15:21:55 Added account UNILEVER with tokens [u un uni unil unile unilev]
2022/11/22 15:21:55 Added account MILLER COORS with tokens [  m mi mil mill mille miller miller  miller c miller co miller coo]
2022/11/22 15:21:55 Listening with HTTP server on :8999
2022/11/22 15:21:58 Adding Transaction &{DANNON %!s(int=300) 2022-11-22 15:21:58.106232236 -0600 CST m=+2.674637854}
2022/11/22 15:21:58 Payer transaction DANNON : 300
2022/11/22 15:21:58 Adding Transaction &{UNILEVER %!s(int=200) 2022-11-22 15:21:58.110808748 -0600 CST m=+2.679214370}
2022/11/22 15:21:58 Payer transaction UNILEVER : 200
2022/11/22 15:21:58 Adding Transaction &{DANNON %!s(int=-200) 2022-11-22 15:21:58.115575194 -0600 CST m=+2.683980818}
2022/11/22 15:21:58 Payer transaction DANNON : -200
2022/11/22 15:21:58 Adding Transaction &{MILLER COORS %!s(int=10000) 2022-11-22 15:21:58.120533682 -0600 CST m=+2.688939306}
2022/11/22 15:21:58 Payer transaction MILLER COORS : 10000
2022/11/22 15:21:58 Adding Transaction &{DANNON %!s(int=1000) 2022-11-22 15:21:58.124458167 -0600 CST m=+2.692863787}
2022/11/22 15:21:58 Payer transaction DANNON : 1000
2022/11/22 15:22:00 Apply 300 points from Payer DANNON and resulting in a points allocation balance of 4700
2022/11/22 15:22:00 Apply 200 points from Payer UNILEVER and resulting in a points allocation balance of 4500
2022/11/22 15:22:00 Apply -200 points from Payer DANNON and resulting in a points allocation balance of 4700
2022/11/22 15:22:00 Apply 4700 points from Payer MILLER COORS and resulting in a points allocation balance of 0
2022/11/22 15:22:00 Payer DANNON being credited 100
2022/11/22 15:22:00 Payer transaction DANNON : -100
2022/11/22 15:22:00 Payer UNILEVER being credited 200
2022/11/22 15:22:00 Payer transaction UNILEVER : -200
2022/11/22 15:22:00 Payer MILLER COORS being credited 4700
2022/11/22 15:22:00 Payer transaction MILLER COORS : -4700
```
