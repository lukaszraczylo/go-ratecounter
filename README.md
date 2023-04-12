## go-ratecounter

[![Run unit tests](https://github.com/lukaszraczylo/go-ratecounter/actions/workflows/test.yaml/badge.svg)](https://github.com/lukaszraczylo/go-ratecounter/actions/workflows/test.yaml) [![Go Reference](https://pkg.go.dev/badge/github.com/lukaszraczylo/go-ratecounter.svg)](https://pkg.go.dev/github.com/lukaszraczylo/go-ratecounter)

The library was inspired heavily by the [ratecounter](https://github.com/paulbellamy/ratecounter) which was unfortunately abandoned for a while. The library is a simple rate counter that can be used to count events over a given time period. The library is thread safe and can be used in a concurrent environment.

### Usage

You can add this library to your project by running the following command:

`go get github.com/lukaszraczylo/go-ratecounter`

### Getting started

#### Single rate counter

The simple configuration whenever you need only one rate counter to exist.
Please note that the default interval is **60 seconds**.

```go
package main

import (
  "fmt"
  "time"

  "github.com/lukaszraczylo/go-ratecounter"
)

func main() {
  // Create a new rate counter that counts over a 60 second period by default
  rc := goratecounter.NewRateCounter()
  // Modify it to count over a 10 second period
  rc.WithConfig(goratecounter.RateCounterConfig{
    Interval: 10 * time.Second,
  })

  // Increment the 'default' counter
  rc.Incr(1)
  // Increment the 'default' counter by 10
  rc.Incr(10)
  // get value of the default counter
  fmt.Println(rc.Get())
  // get number of increments of the default counter
  fmt.Println(rc.GetTicks())
  // get rate of the default counter
  fmt.Println(rc.GetRate())
  // get average increment of the default counter
  fmt.Println(rc.Average())
}
```

#### Multiple rate counters in one

The advanced configuration whenever you need multiple rate counters to exist.

```go
package main

import (
  "fmt"
  "time"

  "github.com/lukaszraczylo/go-ratecounter"
)

func main() {
  // Create a new rate counter that counts over a 60 second period by default
  rc := goratecounter.NewRateCounter()
  // Modify it to count over a 10 second period
  rc.WithConfig(goratecounter.RateCounterConfig{
    Interval: 10 * time.Second,
  })
  rc.WithName("testing-123")
  rc.WithName("testing-456")

  // Increment the the newly created counters
  rc.IncrByName("testing-123", 1)
  rc.IncrByName("testing-456", 7)

  // get value of the default counter
  fmt.Println(rc.Get()) // 0
  // get value of the 'testing-123' counter
  fmt.Println(rc.GetByName("testing-123")) // 1
  // get value of the 'testing-456' counter
  fmt.Println(rc.GetByName("testing-456")) // 7
}
```

