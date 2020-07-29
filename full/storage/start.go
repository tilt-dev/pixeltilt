package main

import (
"fmt"
"time"
)

func Benchmark() {
start := time.Unix(0, 1596050557060185584)
if start.IsZero() {
fmt.Println("Couldn't benchmark start time!")
}
 	fmt.Println("storage service restarted in:", time.Since(start).Round(time.Millisecond))
}
