# Metrics

Demonstation of using metrics reporter.  

This demo should be used among with graphite/grafana.
I used these docker images to spin-up:
- graphiteapp/graphite-statsd
- grafana/grafana


If you find something strange, it's better to investigate the problem this way:  
1. Look at logs of statsd. Also add `"dumpMessages": true` to it's config
2. Look what statsd sends to graphite. I used `tcpdump -i lo -A tcp port 2003`.


## Building
```
$ dep ensure
$ go build
```

## Using
Launch `./metrics`. It will produce many logs to stdout.  
In another terminal window run simulation:
```
while :; do
  curl -s http://localhost:10000/hello2 >/dev/null; sleep 0.2
done
```

## Preparing grafana
- Add graphite source
- Do some settings (theme, etc.)
- Import dashboard from `dashboards/`

### errors panel in grafana
Also there is some errors. You can simulate 'em by requesting `/errors?who=cause`.
  The "cause" may be:
- redis
- mysql
- file_write

# WARNING
Never use Max OS Docker for that. Just [because](https://stackoverflow.com/q/53400513/10682059).