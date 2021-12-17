#!/bin/sh
docker push ruqqq/prayertime-ha-armv7:latest
docker push ruqqq/prayertime-ha-armv7:${jq --raw-output '.version' config.json}
