package main

import (
"fmt"
"time"
)

func Benchmark() {
start := time.Unix(0, 1586355613655333180)
if start.IsZero() {
fmt.Println("Couldn't benchmark start time!")
}
 	fmt.Println("rectangler service restarted in:", time.Since(start).Round(time.Millisecond))
}
