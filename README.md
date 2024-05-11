# Online Watch

Monitors targets for packetloss using TCP and ICMP(v6). Designed to be run inside DCs to track durations of packet-loss during maintenance.

```
Usage: online-watch [-p|--port=<port>] [--no-tcp] [--no-icmp] [--verbose] [--interval=<ms>] [--timeout=<ms>] TARGETS...

Monitors targets for packetloss using TCP and ICMP(v6).

Arguments:
  TARGETS          [label=]<IPv4|IPv6>[:port]

Options:
  -p, --port       Disable TCP checks (default 22)
      --no-tcp     Disable TCP checks
      --no-icmp    Disable ICMP checks
      --verbose    Log all results
      --interval   Time between checks (ms) (default 250)
      --timeout    TCP/ICMP timeout (ms) (default 250)
```

## Output

To watch 10.0.0.1 for TCP/22 and ICMP, only showing output during outages:

```console
$ online-watch 10.0.0.1
2024-05-11 08:55:54.132 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 0s (1x)
2024-05-11 08:55:54.378 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 246ms (2x)
2024-05-11 08:55:54.628 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 497ms (3x)
2024-05-11 08:55:54.882 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 750ms (4x)
2024-05-11 08:55:55.137 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 1.005s (5x)
2024-05-11 08:55:55.382 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 1.251s (6x)
2024-05-11 08:55:55.634 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 1.503s (7x)
2024-05-11 08:55:55.878 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 1.747s (8x)
2024-05-11 08:55:56.136 10.0.0.1 ICMP ✔︎ TCP/22 ✖︎ for 2.004s (9x)
^C
- Longest outage 10.0.0.1: 2.004s (9x)
```

## Installation

```
go install github.com/href/online-watch@latest
```

## Examples

To check multiple IPs concurrently:
```
online-watch 10.0.0.1 10.0.0.2
```

To attach labels to each IP for more helpful output:
```
online-watch rack-a=10.0.0.1 rack-b=10.0.0.2
```

To use a different default port:
```
online-watch -p 80 web-a=10.0.0.1 web-b=10.0.0.2
```

To use a separate port per target:
```
online-watch frontend=10.0.0.1:443 backend=10.0.0.2:5432
```
