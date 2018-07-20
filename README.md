tcpproxy
====================
tcp proxy in golang

## usage
```shell
Usage: tcpproxy --portmaps PORTMAPS [--limit LIMIT] [--keepalive] [--aliveperiod ALIVEPERIOD]

Options:
  --portmaps PORTMAPS, -m PORTMAPS
  --limit LIMIT, -l LIMIT
  --keepalive, -k
  --aliveperiod ALIVEPERIOD, -a ALIVEPERIOD [default: 180ns]
  --help, -h             display this help and exit
```


for example:

```bash
tcpproxy -m 127.0.0.1:8228:192.168.1.39:8228 -m 0:8229:192.168.1.39:8229
```