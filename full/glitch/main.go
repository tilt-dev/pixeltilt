package main

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sug0/go-glitch"
	"github.com/windmilleng/enhance/render/api"
)

// glitch formulas from: https://github.com/sug0/go-glitch/blob/master/res/cool.txt

// const glitchExprStr = "128 & (c - ((c - 150 + s) > 5 < s))"
// const glitchExprStr = "(c & (c ^ 55)) + 25"
// const glitchExprStr = "r | ((255 - (Y > 2) : ((r | c) < 2)) - 140)"
// const glitchExprStr = "(Y | (c > 1)) ^ 128"
// const glitchExprStr = "e ^ c - (e - 55)"
// const glitchExprStr = "86 ^ ((R&c) > ((G&c) - ((G&c) ^ ((G&c) / (175 : (x < ((B&c) + ((G&c) + r))))))))"
// const glitchExprStr = "(Y - ((G & 55)|(R - 25))) ^ 25"
// const glitchExprStr = "Y : (146 - (e | (15 ? 185)))"
// const glitchExprStr = "Y & (206 / (e & (Y / r)))"
// const glitchExprStr = "179 ^ (Y % (x | (s : (209 : (Y % (G - (r # r)))))))"
// const glitchExprStr = "36 - (67 | (Y - (G | R)))"
// const glitchExprStr = "128 & (c - (s - 255) + s) : (s ^ (c ^ 255)) + 25"
// const glitchExprStr = "((Y/s) + c < 1) + (x*x*y) % (c ^ ((y*x - y*y) < (y > 5)))"
// const glitchExprStr = "(c&x-y)/(s | x*y) - ((y+c)/x < 2)"
// const glitchExprStr = "(c & (s ^ 55)) + (25 > s)"
// const glitchExprStr = "128 ^ ((r ^ 15) | (c - (s - (r * 255))))"
// const glitchExprStr = "(r | c) < 3"
// const glitchExprStr = "(r - Y) | (((r | ((c | Y) > 2)) < 2) + (Y % 255))"
// const glitchExprStr = "((Y | ((((108 - Y) * s) > c + Y + s) ^ 5 + 25)) * (r > 5)) : ((c < 1) | Y)"
// const glitchExprStr = "(s & c) | (((Y + (r - 55)) ^ s) > 10) < 10 - ((128 - c) | Y)"
// const glitchExprStr = "(s & c) | (((Y + (r - 55)) ^ s) > 15) < 10 - ((16 - Y) | c)"
// const glitchExprStr = "(r - ((Y % 15) : (r < 5)) > 10) | (s - c)"
// const glitchExprStr = "c|((e - s + x - y) / N) - (R&c)|(B&c) + ((G&c)*(G&c)*(G&c))/(x*x*x) - (s & (R&c))"
// const glitchExprStr = "((x*x*x % y) + (c - s) + r) > 1"
// const glitchExprStr = "((x - y) # (y - x) ? (R&c)|(G&c)) - ((255 - (80 ? c)) & c)"
// const glitchExprStr = "((128 ? x) ^ (128 ? y)) : ((y - x) # (x - y) ? (Y-c)) - ((255 - (80 ? s)) & y - (R&c)|(G&c))"
// const glitchExprStr = "((s - x) # (e - s)) ? ((128 - x) ^ (y - c)) ? (128 - Y)"
// const glitchExprStr = "c & (((230 - (130 ? Y)) : (s - ((G&c)|(R&c))) + s) - s) + s"
// const glitchExprStr = "y ^ ((G&c) + ((B&c) | ((B&c) * (y | (166 ? ((G&c) % s))))))"
// const glitchExprStr = "33 + (e - (s > ((R&c) : (Y + (x ? (204 < 243))))))"
// const glitchExprStr = "(e | R) - ((c & G) - (Y & G))"
// const glitchExprStr = "b ^ (233 < (B > (r # b)))"
// const glitchExprStr = "Y ^ (106 & (Y ? (r < G)))"
// const glitchExprStr = "182 + (x + (s - (y % (15 & (r ? (e + (s : (76 - (y # (r * (r ? (r | (y % (r ? (195 ? (R - (123 > (b : N)))))))))))))))")))
// const glitchExprStr = "b ^ (164 - (b < (G # (b ^ (e % r)))))"
// const glitchExprStr = "c | (s < (c & (r ? (B | (e - (Y < (Y ^ Y)))))))"
// const glitchExprStr = "b ^ (r | (s : (x # B)))"
// const glitchExprStr = "G / (b / (N > (110 > s)))"
// const glitchExprStr = "29 ? (Y & (c ^ (s > c)))"
// const glitchExprStr = "((c|Y) < 3) @ ((H-L+10)|(((c < 1) @ ((s/x) * y)) ^ (c - L) | (Y - L)))"
// const glitchExprStr = "(c<(s%4)>1)@(L+25)-25"
// const glitchExprStr = "(c&(R-Y)&25)^(Y|((Y:h)-25))@(Y-s)"
// const glitchExprStr = "255 ^ ((H ^ L) - s) ? ((c & (R @ 246)) | (c & (G @ 155)) | (c & (B @ 255)))"
// const glitchExprStr = "e@((H-((c < 2) - (N@r)%128))+(L&R)-(H-G))-((e&r)/16)"
// const glitchExprStr = "((x*x*x*y*y*y) < s) + (c > 4)"
// const glitchExprStr = "c & (Y - x*y)"
// const glitchExprStr = "c * (x ^ 5)"
// const glitchExprStr = "((y*s - y*s) & (Y | c)) ^ (x|c)"
// const glitchExprStr = "((y*s - y*s < s) & (Y | c)) ^ (x|c) > s"
// const glitchExprStr = "(r - y) | ((55 - x*3) ^ 25) : (((y*s - y*s < s) & (Y | c)) ^ ((255-y)+x|c) > s)"
// const glitchExprStr = "128 + ((r > 2 < 1) : ((x|y) ^ 55)) % 50 - c + y"
// const glitchExprStr = "c & ((x+y*y*y+r) : (s + c % 15) - x*x)"
// const glitchExprStr = "c ^ ((x ? Y) ? (255 - e))"
// const glitchExprStr = "((Y - x) ? (c - s)) ^ (0 - (x + y) : (0 - (y - x)))"
// const glitchExprStr = "((s + (R&c)) - c) & (0 - (x ? (0 - (x - x*y) : y)) - (1 : (y ^ x))) - (c + s)"
// const glitchExprStr = "(c @ (s & G)) ^ (c - L)"
// const glitchExprStr = "L ^ H > H < L % R"
// const glitchExprStr = "(c ^ (L ^ H)) % (R > c) : (( L ^ H < R > c) % (5 ^ R))"
// const glitchExprStr = "(Y - c)|(r - s/2) + (61 & (R - (r % Y)))"
// const glitchExprStr = "(c-(s>2))/16*((Y-H)@128)"
// const glitchExprStr = "((y/(x*x) ^ s) + x*s) > 2 + c"
// const glitchExprStr = "Y < (s # (e * (y & (x : (x & (s % (r + (s ? c))))))))"
// const glitchExprStr = "120 : (Y & (N | (s & (B & (e # (Y + (Y : R)))))))"
// const glitchExprStr = "Y | (34 % (194 < (e < G)))"
// const glitchExprStr = "r + (61 & (R - (r % e)))"
// const glitchExprStr = "Y ^ (s & (60 > (y < (239 ? (c / b)))))"
// const glitchExprStr = "#((1 - y : x) ^ (1 - y & x)) & (0 - (x*x*x % y))"
// const glitchExprStr = "(s ^ (1 - x : y)) & ((((e|(G&c)) - r)|N < 2) > 1)"
// const glitchExprStr = "((1 - y : x) | (1 - x : y)) - c"
// const glitchExprStr = "(R&c)|(G&c) % (((1 - y : x) ^ (1 - x : y)) - c)"
// const glitchExprStr = "100 + ((N & 25) ^ (1 - x : y)) + (75 & Y)"
// const glitchExprStr = "((G&c) ^ (e - (R&c))) % ((N & 25) ^ (e - x : y)) + (75 & Y)"
// const glitchExprStr = "(Y : c) ^ ((c + 50) & (0 - (y ? (0 - x : y)) & (y : x)))"
const glitchExprStr = "(255 - (140 ? c)) & c"

// const glitchExprStr = "(H-L)|b"
// const glitchExprStr = "L ^ H"
// const glitchExprStr = "(c ^ (L ^ H)) % (R > c) : ( L ^ H < R > c)"
// const glitchExprStr = "(r > (H ^ L) % b) : (H % L)"
// const glitchExprStr = "((H + (s/L) - L) & R) | ((G / (b / (N > (110 > s)))) : R)"
// const glitchExprStr = "b ^ (r | (s : (x # B)))"

var glitchExpr *glitch.Expression

func main() {
	Benchmark()
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var err error
	glitchExpr, err = glitch.CompileExpression(glitchExprStr)
	if err != nil {
		log.Fatal(err)
	}

	HandleWithNoSubpath("/", api.HttpRenderHandler(render))
	http.HandleFunc("/set_expr", setExpr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func HandleWithNoSubpath(path string, f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(w, req)
			return
		}
		f(w, req)
	}

	http.HandleFunc(path, handler)
}

// just a dumb form to let people play around with different expressions without restarting
func setExpr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(`
<form method="post">
expr: <input type="text" name="expr">
<input type="submit" value="Submit">
</form>
`))
	if err != nil {
		log.Printf("error writing response: %v\n", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing request: %v", err), http.StatusBadRequest)
		return
	}

	expr := r.PostForm.Get("expr")
	if expr != "" {
		g, err := glitch.CompileExpression(expr)
		if err != nil {
			_, err := w.Write([]byte(fmt.Sprintf("<br>Invalid expression: %s<br>", err.Error())))
			if err != nil {
				log.Printf("error writing response: %v\n", err)
				return
			}
		}
		glitchExpr = g
		log.Printf("changed expression to %s\n", expr)
		_, err = w.Write([]byte(fmt.Sprintf("<br>Expression set to %s<br>", expr)))
		if err != nil {
			log.Printf("error writing response: %v\n", err)
			return
		}
	}
}

func render(req api.RenderRequest) (api.RenderReply, error) {
	input, err := png.Decode(bytes.NewReader(req.Image))
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "decoding image")
	}

	output, err := glitchExpr.JumblePixels(input)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "glitching image")
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, output)
	if err != nil {
		return api.RenderReply{}, errors.Wrap(err, "encoding image")
	}

	return api.RenderReply{Image: buf.Bytes()}, nil
}
