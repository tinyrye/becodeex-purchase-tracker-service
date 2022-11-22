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

### Building ###

Execute `build.sh`.  This builds a `run_http_service` executable.