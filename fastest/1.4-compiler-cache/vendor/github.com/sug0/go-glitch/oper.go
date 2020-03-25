package glitch

type operator struct {
    precedence int
    assoc      int
    f          func(uint8, uint8) uint8
}

const (
    assocRight = iota
    assocLeft
)

var operMap = map[rune]*operator{
    '+': &operator{precedence: 4, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x + y
    }},
    '-': &operator{precedence: 4, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x - y
    }},
    '|': &operator{precedence: 4, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x | y
    }},
    '^': &operator{precedence: 4, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x ^ y
    }},
    '*': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x * y
    }},
    '/': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        if y == 0 {
            return x
        } else {
            return x / y
        }
    }},
    '%': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        if y == 0 {
            return x
        } else {
            return x % y
        }
    }},
    '<': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x << y
    }},
    '>': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x >> y
    }},
    '&': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x & y
    }},
    ':': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        return x &^ y
    }},
    '#': &operator{precedence: 5, assoc: assocLeft, f: func(x, y uint8) uint8 {
        var z uint8 = 1
        for y > 0 {
            z *= x
            y--
        }
        return z
    }},
    '?': &operator{precedence: 6, assoc: assocRight, f: func(x, y uint8) uint8 {
        if x > y {
            return 255
        } else {
            return 0
        }
    }},
    '@': &operator{precedence: 6, assoc: assocRight, f: func(x, y uint8) uint8 {
        fuzz := float64(y)/255.0
        return uint8(float64(x) * fuzz)
    }},
}

func (o1 *operator) hasPrecedence(o2 *operator) bool {
    return (o2.assoc == assocRight && o1.precedence > o2.precedence) ||
           (o2.assoc == assocLeft && o1.precedence >= o2.precedence)
}
