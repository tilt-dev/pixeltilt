package main

import (
"fmt"
"time"
)

func Benchmark() {
start := time.Unix(0, 1626887492364184905)
if start.IsZero() {
fmt.Println("Couldn't benchmark start time!")
}
 	fmt.Println("render service restarted in:", time.Since(start).Round(time.Millisecond))
}
