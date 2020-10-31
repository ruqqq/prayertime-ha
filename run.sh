#!/usr/bin/with-contenv bashio
CONFIG_PATH=/data/options.json

export REMINDER_OFFSET="$(jq --raw-output '.reminder_offset' $CONFIG_PATH)"
export LOCATION_CODE="$(jq --raw-output '.location_code' $CONFIG_PATH)"

./prayertime_ha

