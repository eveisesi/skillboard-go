#!/bin/sh

cp /data/tmpl_redis.conf /etc/redis.conf

sed -i "s/REPLACE_ME_REDIS_PASS/$REDIS_PASS/g" /etc/redis.conf

redis-server /etc/redis.conf