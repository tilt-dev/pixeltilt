package glitch

import (
    "image"
    "image/gif"
    "image/color"
    "image/draw"
)

func (expr *Expression) JumblePixels(data image.Image) (image.Image, error) {
    return expr.JumblePixelsMonitor(data, nil)
}

func (expr *Expression) JumblePixelsMonitor(data image.Image, mon func(int, int)) (image.Image, error) {
    var i, pixsize int

    if mon != nil {
        pixsize = data.Bounds().Dx() * data.Bounds().Dy()
    }

    return expr.jumblePixels(data, &i, pixsize, mon)
}

// for now support only gifs with the same dimensions per frame
func (expr *Expression) JumbleGIFPixels(data *gif.GIF) (*gif.GIF, error) {
    return expr.JumbleGIFPixelsMonitor(data, nil)
}

func (expr *Expression) JumbleGIFPixelsMonitor(data *gif.GIF, mon func(int, int)) (*gif.GIF, error) {
    var i, pixsize, imgsize int

    if mon != nil {
        imgsize = data.Image[0].Bounds().Dx() * data.Image[0].Bounds().Dy()
        pixsize = imgsize * len(data.Image)
    }

    newgif := new(gif.GIF)
    newgif.LoopCount = data.LoopCount
    newgif.Delay = make([]int, len(data.Delay))
    copy(newgif.Delay, data.Delay)

    overpaint := image.NewNRGBA(data.Image[0].Bounds())
    var didFirst bool

    for _,img := range data.Image {
        // glitch image
        newimg, err := expr.jumblePixels(img, &i, pixsize, mon)
        if err != nil {
            return nil, err
        }

        if didFirst {
            draw.Draw(overpaint, overpaint.Bounds(), newimg, image.ZP, draw.Over)
        } else {
            didFirst = true
            draw.Draw(overpaint, overpaint.Bounds(), newimg, image.ZP, draw.Src)
        }
        newimg = nil

        // convert from NRGBA to Paletted
        paletted, err := imgNRGBAToPaletted(overpaint)
        if err != nil {
            return nil, err
        }
        newgif.Image = append(newgif.Image, paletted)
    }

    return newgif, nil
}

func (expr *Expression) jumblePixels(data image.Image, i *int, pixsize int, mon func(int, int)) (image.Image, error) {
    img := image.NewNRGBA(data.Bounds())
    xm, xM := data.Bounds().Min.X, data.Bounds().Max.X
    ym, yM := data.Bounds().Min.Y, data.Bounds().Max.Y

    var err error
    var nr, ng, nb, sr, sg, sb uint8

    if mon != nil {
        mon(*i, pixsize)
    }

    for x := xm; x < xM; x++ {
        for y := ym; y < yM; y++ {
            r, g, b, a := convUint8(data.At(x, y).RGBA())
            nr, ng, nb, err = expr.evalRPN(x, y, xM, yM,
                                           r, g, b, a,
                                           sr, sg, sb, data)

            if err != nil {
                return nil, err
            }

            sr = nr
            sg = ng
            sb = nb

            img.Set(x, y, color.NRGBA{
                nr, ng, nb, a,
            })

            if mon != nil {
                mon(*i, pixsize)
                *i++
            }
        }
    }

    return image.Image(img), nil
}
