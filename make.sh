#!/bin/sh
GOOS=linux GOARCH=arm GOARM=7 go build -o prayertime_ha

docker run --rm -ti --name hassio-builder --privileged \
  -v .:/data \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  homeassistant/amd64-builder -t /data --all --test \
  -i prayertime-ha-{arch} -d ruqqq
