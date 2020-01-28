package glitch

import (
    "bytes"
    "image"
    "image/gif"
    "sync"
)

var bufPool = sync.Pool{
    New: func() interface{} { return &bytes.Buffer{} },
}

// pretty shit hack to get a Paletted image out of an NRGBA image
func imgNRGBAToPaletted(data *image.NRGBA) (*image.Paletted, error) {
    buf := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(buf)

    if err := gif.Encode(buf, data, nil); err != nil {
        return nil, err
    }

    decoded, err := gif.Decode(buf)
    if err != nil {
        return nil, err
    }

    return decoded.(*image.Paletted), nil
}

func convUint8(r, g, b, a uint32) (uint8, uint8, uint8, uint8) {
    return uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), uint8(a / 0x101)
}

func threeRule(x, max int) uint8 {
    return uint8(((255 * x) / max) & 255)
}

func max(vals ...uint8) (m uint8) {
    for _,v := range vals {
        if v > m {
            m = v
        }
    }
    return
}

func min(vals ...uint8) (m uint8) {
    m = 255
    for _,v := range vals {
        if v < m {
            m = v
        }
    }
    return
}

func fetchBox(box *[9]sum, x, y int, r, g, b uint8, data image.Image) {
    k := 0
    for i := x - 1; i <= x + 1; i++ {
        for j := y - 1; j <= y + 1; j++ {
            if i == x && j == y {
                (*box)[k] = sum{r, g, b}
                k++
                continue
            }
            r0, g0, b0,_ := convUint8(data.At(i, j).RGBA())
            (*box)[k] = sum{r0, g0, b0}
            k++
        }
    }
    return
}
