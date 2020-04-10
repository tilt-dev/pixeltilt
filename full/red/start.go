package main

import (
"fmt"
"time"
)

func Benchmark() {
start := time.Unix(0, 1586536031812054010)
if start.IsZero() {
fmt.Println("Couldn't benchmark start time!")
}
 	fmt.Println("red service restarted in:", time.Since(start).Round(time.Millisecond))
}
