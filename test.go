package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
}

func main() {
	fmt.Println("Run server")
	fasthttp.ListenAndServe("8080", fastHTTPHandler)
}
