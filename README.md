# brokenProxy

brokenProxy is a Go application that acts as a proxy server to a single upstream, but exposes that upstream on different URLs. This is useful for testing failover logic, where you want to simulate a cluster of unstable endpoints.

The application listens on a configurable port and periodically sets between one and four simulated hosts to be down, chosen randomly. All requests are forwarded to a single upstream server URL, which is also configurable.

The brokenProxy application exposes eight endpoints, each of which is a simulated host that is assigned a path:

```
/one
/two
/three
/four
/five
/six
/seven
/eight
```

Of these eight endpoints, a random number between one and four of them are set to be down at any given time. Which endpoints are down changes randomly every minute.


## Usage

### Flags

You can set the upstream server address and the local server port via the following command-line flags:

```
    -upstream-addr: Upstream server address
    -port: Local server port (default: 8080)
```

For example:

```
$ brokenProxy -upstream-addr http://localhost:8000 -port 8080

```


### Environment variables

You can also set the upstream server address and the local server port via the following environment variables:

```
    UPSTREAM_ADDR: Upstream server address
    PORT: Local server port (default: 8080)
```

For example:

shell

```
$ export UPSTREAM_ADDR=http://localhost:8000
$ export PORT=8080
$ brokenProxy
```



