package main

import (
"fmt"
"time"
)

func Benchmark() {
start := time.Unix(0, 1586544976718630251)
if start.IsZero() {
fmt.Println("Couldn't benchmark start time!")
}
 	fmt.Println("muxer service restarted in:", time.Since(start).Round(time.Millisecond))
}
