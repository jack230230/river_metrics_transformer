#river_metrics_transformer

transform go-mysql-elasticsearch stat page to prometheus metrics page

## Usage
```
godep go run river_metrics_transformer -host 0.0.0.0 -port 8080 -river http://192.168.0.3:12800/stat
```
Then prometheus can access `0.0.0.0:8080/metrics` to fetch metrics.

_Note_: `host`, `port`, `river` can also be set by environment variables. The environment variable names are `CONFIG_HOST`, `CONFIG_PORT`, `CONFIG_RIVER`

## Docker
There is already a docker image for use:
```
docker pull eaglechen/river_metrics_transformer
```
