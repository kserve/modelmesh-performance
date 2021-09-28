# Trunks [![Build Status](https://travis-ci.org/straightdave/trunks.svg?branch=master)](https://travis-ci.org/straightdave/trunks)

Trunks, like every son, is derived from the father Vegeta with some enhanced skills:
1. dump HTTP reponses
2. gRPC support

![Trunks](http://images2.wikia.nocookie.net/__cb20100725123520/dragonballfanon/images/5/52/Future_Trunks_SSJ2.jpg)

> Current Tag v13.0.0

## Usage manual

for original usage of Vegeta, please refer to [vegeta' readme](https://github.com/tsenart/vegeta/blob/master/README.md)

## More functionalities

### dump http attck response to file
```console
(add one more option '-respf' to 'attack')
-respf string
      Dump responses to file
```

### gRPC perf test (as lib)

Trunks uses the _burn_ to attack the gRPC services. Capable with any gRPC service.

Given multiple hosts ('IP:port' of service instances), Trunks will use simple round-robin as client-side load balance mechanism.

>Burning duration is shorter. No need to do the real __service discovery__ or __watch / real-time live connection's update__. Say, if testing service instances registered to Etcd, we don't need to watch the status of all instances and adjust the connection pool. During the short time of perf testing, we assume the instance hosts are _not-changing_.

Example:

>This example is using "google.golang.org/grpc/examples/route_guide" as the target server.

```golang
package main

import (
    "fmt"
    "time"

    trunks "github.com/straightdave/trunks/lib"
    // for convenience, change the client_stub.pb.go into package main
)

func main() {
    tgt := &trunks.Gtarget{
        MethodName: "/routeguide.RouteGuide/GetFeature",
        Requests:   []proto.Message{&Point{Latitude: 10000, Longitude: 10000}}, // supporting multiple requests
        Response:   &Feature{}, // providing one response data struct
    }

    burner, err := trunks.NewBurner(
        []string{"192.168.0.1:8087"},  // server address pool; simple round-robin
        trunks.WithLooping(true),      // loop requests; false by default
        trunks.WithNumWorker(20),      // worker goroutine pool size; 10 is default
        trunks.WithDumpFile("a.dump"), // dump responses to file
    )
    if err != nil {
        fmt.Println(err)
        return
    }
    defer burner.Close()

    var metrics trunks.Metrics
    startT := time.Now()

    // burning the target service with QPS=5 and Duration=10s
    for res := range burner.Burn(tgt, uint64(5), 10*time.Second) {
        metrics.Add(res)
    }
    dur := time.Since(startT)
    metrics.Close()

    fmt.Printf("dur: %v\n", dur.Seconds())
    fmt.Printf("earliest: %v\n", metrics.Earliest.Sub(startT).Nanoseconds())
    fmt.Printf("latest: %v\n", metrics.Latest.Sub(startT).Nanoseconds())
    fmt.Printf("end: %v\n", metrics.End.Sub(startT).Nanoseconds())
    fmt.Printf("reqs: %d\n", metrics.Requests)
    fmt.Printf("success%: %f\n", metrics.Success)
    fmt.Printf("p50: %s\n", metrics.Latencies.P50)
    fmt.Printf("p95: %s\n", metrics.Latencies.P95)
    fmt.Printf("p99: %s\n", metrics.Latencies.P99)
    fmt.Printf("mean: %s\n", metrics.Latencies.Mean)
    fmt.Printf("max: %s\n", metrics.Latencies.Max)
    // ...
}
```

For this code snippet, it would result in:
```console
dur: 9.802099215
earliest: 68670
latest: 9800068577
end: 9802058490
reqs: 50
success%: 99.88
p50: 5.974748ms
p95: 6.084433ms
p99: 6.10946ms
mean: 5.143272ms
max: 6.19225ms
```

## Arion as gRPC
Arion makes it easy to use Trunks. Please check https://github.com/straightdave/arion


