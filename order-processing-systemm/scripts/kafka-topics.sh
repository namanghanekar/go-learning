#!/bin/bash

docker exec kafka kafka-topics \
--create \
--topic order.created \
--bootstrap-server localhost:9092 \
--partitions 1 \
--replication-factor 1