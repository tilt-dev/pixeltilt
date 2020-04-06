package glitch

import (
    "fmt"
    "math/rand"
    "image"
    "strconv"
)

type sum struct {
    r, g, b uint8
}

// save some values that
// might be expensive to recalculate
type sumsave struct {
    v_Y,
    v_e,
    v_b,
    v_r,
    v_h,
    v_v,
    v_d,
    v_H,
    v_L *sum
}

func (expr *Expression) evalRPN(x, y, w, h int,
                                r, g, b, a uint8,
                                sr, sg, sb uint8,
                                data image.Image,
) (rr uint8, gr uint8, br uint8, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("error evaluating expression: %s", expr.infix)
            return
        }
    }()

    if a == 0 {
        return
    }

    stk := make([]sum, 0, len(expr.toks))
    saved := sumsave{}

    var box [9]sum

    for _,tok := range expr.toks {
        if tok == "" {
            continue
        }
        if oper := operMap[[]rune(tok)[0]]; oper != nil {
            a0, b0 := stk[len(stk)-2], stk[len(stk)-1]
            stk = append(stk[:len(stk)-2], sum{oper.f(a0.r, b0.r),
                                               oper.f(a0.g, b0.g),
                                               oper.f(a0.b, b0.b)})
        } else if tok == "c" {
            stk = append(stk, sum{r, g, b})
        } else if tok == "R" {
            stk = append(stk, sum{255, 0, 0})
        } else if tok == "G" {
            stk = append(stk, sum{0, 255, 0})
        } else if tok == "B" {
            stk = append(stk, sum{0, 0, 255})
        } else if tok == "Y" {
            if saved.v_Y == nil {
                y := uint8(float64(r)*0.299) +
                     uint8(float64(g)*0.587) +
                     uint8(float64(b)*0.0722)
                saved.v_Y = &sum{y, y, y}
            }
            stk = append(stk, *saved.v_Y)
        } else if tok == "s" {
            stk = append(stk, sum{sr, sg, sb})
        } else if tok == "x" {
            xu := threeRule(x, w)
            stk = append(stk, sum{xu, xu, xu})
        } else if tok == "y" {
            yu := threeRule(y, h)
            stk = append(stk, sum{yu, yu, yu})
        } else if tok == "r" {
            if saved.v_r == nil {
                i0, j0 := rand.Int() % 3 - 1, rand.Int() % 3 - 1
                i1, j1 := rand.Int() % 3 - 1, rand.Int() % 3 - 1
                i2, j2 := rand.Int() % 3 - 1, rand.Int() % 3 - 1
                r0,_,_,_ := convUint8(data.At(x + i0, y + j0).RGBA())
                _,g0,_,_ := convUint8(data.At(x + i1, y + j1).RGBA())
                _,_,b0,_ := convUint8(data.At(x + i2, y + j2).RGBA())
                saved.v_r = &sum{r0, g0, b0}
            }
            stk = append(stk, *saved.v_r)
        } else if tok == "e" {
            if saved.v_e == nil {
                fetchBox(&box, x, y, r, g, b, data)
                dr := box[8].r - box[0].r + box[5].r - box[3].r +
                      box[7].r - box[1].r + box[6].r - box[2].r
                dg := box[8].g - box[0].g + box[5].g - box[3].g +
                      box[7].g - box[1].g + box[6].g - box[2].g
                db := box[8].b - box[0].b + box[5].b - box[3].b +
                      box[7].b - box[1].b + box[6].b - box[2].b
                saved.v_e = &sum{dr, dg, db}
            }
            stk = append(stk, *saved.v_e)
        } else if tok == "b" {
            if saved.v_b == nil {
                fetchBox(&box, x, y, r, g, b, data)
                sr := int(box[0].r) + int(box[1].r) + int(box[2].r) + int(box[3].r) +
                      int(box[4].r) + int(box[5].r) + int(box[6].r) + int(box[7].r) + int(box[8].r)
                sg := int(box[0].g) + int(box[1].g) + int(box[2].g) + int(box[3].g) +
                      int(box[4].g) + int(box[5].g) + int(box[6].g) + int(box[7].g) + int(box[8].g)
                sb := int(box[0].b) + int(box[1].b) + int(box[2].b) + int(box[3].b) +
                      int(box[4].b) + int(box[5].b) + int(box[6].b) + int(box[7].b) + int(box[8].b)
                saved.v_b = &sum{uint8(sr/9), uint8(sg/9), uint8(sb/9)}
            }
            stk = append(stk, *saved.v_b)
        } else if tok == "H" {
            if saved.v_H == nil {
                fetchBox(&box, x, y, r, g, b, data)
                rM := max(box[0].r, box[1].r, box[2].r, box[3].r,
                          box[4].r, box[5].r, box[6].r, box[7].r, box[8].r)
                gM := max(box[0].g, box[1].g, box[2].g, box[3].g,
                          box[4].g, box[5].g, box[6].g, box[7].g, box[8].g)
                bM := max(box[0].b, box[1].b, box[2].b, box[3].b,
                          box[4].b, box[5].b, box[6].b, box[7].b, box[8].b)
                saved.v_H = &sum{rM, gM, bM}
            }
            stk = append(stk, *saved.v_H)
        } else if tok == "L" {
            if saved.v_L == nil {
                fetchBox(&box, x, y, r, g, b, data)
                rM := min(box[0].r, box[1].r, box[2].r, box[3].r,
                          box[4].r, box[5].r, box[6].r, box[7].r, box[8].r)
                gM := min(box[0].g, box[1].g, box[2].g, box[3].g,
                          box[4].g, box[5].g, box[6].g, box[7].g, box[8].g)
                bM := min(box[0].b, box[1].b, box[2].b, box[3].b,
                          box[4].b, box[5].b, box[6].b, box[7].b, box[8].b)
                saved.v_L = &sum{rM, gM, bM}
            }
            stk = append(stk, *saved.v_L)
        } else if tok == "N" {
            rn, gn, bn := uint8(rand.Int() % 256), uint8(rand.Int() % 256), uint8(rand.Int() % 256)
            stk = append(stk, sum{rn, gn, bn})
        } else if tok == "h" {
            if saved.v_h == nil {
                rh, gh, bh,_ := convUint8(data.At(w - x - 1, y).RGBA())
                saved.v_h = &sum{rh, gh, bh}
            }
            stk = append(stk, *saved.v_h)
        } else if tok == "v" {
            if saved.v_v == nil {
                rv, gv, bv,_ := convUint8(data.At(x, h - y - 1).RGBA())
                saved.v_v = &sum{rv, gv, bv}
            }
            stk = append(stk, *saved.v_v)
        } else if tok == "d" {
            if saved.v_d == nil {
                rd, gd, bd,_ := convUint8(data.At(w - x - 1, h - y - 1).RGBA())
                saved.v_d = &sum{rd, gd, bd}
            }
            stk = append(stk, *saved.v_d)
        } else {
            if i, err := strconv.Atoi(tok); err == nil {
                u := uint8(i)
                stk = append(stk, sum{u, u, u})
            }
        }
    }

    v := stk[len(stk)-1]
    stk = nil

    return v.r, v.g, v.b, nil
}
