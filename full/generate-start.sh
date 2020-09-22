#!/bin/bash

for val in $services; do
read -r -d '' VAR <<- EOM
package main

import (
	"fmt"
	"time"
)

func Benchmark() {
	start := time.Unix(0, $(gdate +%s%N))
	if start.IsZero() {
		fmt.Println("Couldn't benchmark start time!")
	}
 	fmt.Println("$(echo $val) service restarted in:", time.Since(start).Round(time.Millisecond))
}
EOM
echo "$VAR" > $(echo "$val""/start.go")
done
