![](https://u.sicp.me/k4Tum.png)

# Go Glitch

A go package to glitch an image based on an expression you pass in as input!
The only limit is your creativity (and patience)!

## Discord

[Join our discord server!](https://discord.gg/cUW4b6Z)

    
# What is the deal with the expressions?

You can think of the image as a functor that you map an expression to, for each pixel's component colors,
returning a new one. The allowed operators are:

* `+` plus
* `-` minus
* `*` multiplication
* `/` division
* `%` modulo
* `#` power of operator
* `&` bit and
* `|` bit or
* `:` bit and not
* `^` bit xor
* `<` bit left shift
* `>` bit right shift
* `?` returns 255 if left side is greater otherwise 0
* `@` attributes a weight in the range `[0, 255]` to the value on the left

The expressions are made up of operators, numbers, parenthesis, and a set of parameters:

* `c` the current value of each pixel component color
* `b` the blurred version of `c`
* `h` the horizontally flipped version of `c`
* `v` the vertically flipped version of `c`
* `d` the diagonally flipped version of `c`
* `Y` the luminosity, or grayscale component of each pixel
* `N` a noise pixel (i.e. a pixel where each component is a random value)
* `R` the red color (i.e. rgb(255, 0, 0))
* `G` the green color (i.e. rgb(0, 255, 0))
* `B` the blue color (i.e. rgb(0, 0, 255))
* `s` the value of each pixel's last saved evaluated expression
* `r` a pixel made up of a random color component from the neighboring 8 pixels
* `e` the difference of all pixels in a box, creating an edge-like effect
* `x` the current x coordinate being evaluated normalized in the range `[0, 255]`
* `y` the current y coordinate being evaluated normalized in the range `[0, 255]`
* `H` the highest valued color component in the neighboring 8 pixels
* `L` the lowest valued color component in the neighboring 8 pixels

## Examples

* `128 & (c - ((c - 150 + s) > 5 < s))`
* `(c & (c ^ 55)) + 25`
* `128 & (c + 255) : (s ^ (c ^ 255)) + 25`

More examples can be found
[here](https://github.com/sugoiuguu/go-glitch/blob/master/res/cool.txt).


# Command line tool usage

Use the [docker image](https://hub.docker.com/r/sugoiuguu/go-glitch/):

    $ docker run sugoiuguu/go-glitch


# API usage

```go
package main

import (
    "os"
    "time"
    "math/rand"
    "image"
    "image/png"
    _ "image/jpeg"

    "github.com/sugoiuguu/go-glitch"
)

func main() {
    f, err := os.Open(os.Args[2])
    if err != nil {
        panic(err)
    }
    defer f.Close()

    f2, err := os.Create(os.Args[1])
    if err != nil {
        panic(err)
    }
    defer f2.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    rand.Seed(time.Now().UnixNano())

    expr, err := glitch.CompileExpression("(n|(c > 1)) ^ 128")
    if err != nil {
        panic(err)
    }

    g, err := expr.JumblePixels(img)
    if err != nil {
        panic(err)
    }
    png.Encode(f2, g)
}
```

# C/C++ interface

There is an experimental interface for C/C++ code
[here](https://github.com/sugoiuguu/go-glitch/blob/master/example/ffi.go).
The API is:

```c
typedef struct {
    char *data;
    size_t size;
} Image_t;

extern Image_t *jumble_pixels(char *expression, char *data, int size);
```
