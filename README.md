# PrayerTime HA Addon
Emits 'prayertime' event on prayer time and on configurable reminder offset.

Shape of payload:
```yaml
prayer_id: 1
prayer_name: Syuruk
is_reminder: false
```

## Developer
### Build
1. Install golang
2. Run `make.sh`
3. Copy `config.json`, `Dockerfile`, `run.sh`, `prayertime_ha` to `/addons/prayertime-ha`
4. Load addon and build
