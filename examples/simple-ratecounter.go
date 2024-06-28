package main

import (
	goratecounter "github.com/lukaszraczylo/go-ratecounter"
)

func main() {
	rc := goratecounter.NewRateCounter()
	rc.Incr(1)
	rc.Get()
}
