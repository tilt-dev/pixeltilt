package main

import (
"fmt"
"time"
)

func Benchmark() {
start := time.Unix(0, 1586445828773995061)
if start.IsZero() {
fmt.Println("Couldn't benchmark start time!")
}
 	fmt.Println("storage service restarted in:", time.Since(start).Round(time.Millisecond))
}
