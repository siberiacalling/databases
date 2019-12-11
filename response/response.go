package response

import (
	"encoding/json"

	"net/http"

	"log"

	"github.com/valyala/fasthttp"
)

func Write(ctx *fasthttp.RequestCtx, statusCode int, body interface{}) {
	ctx.SetContentType("application/json")

	resp, err := json.Marshal(body)
	if err != nil {
		log.Println(err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if _, err := ctx.Write(resp); err != nil {
		log.Println(err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(statusCode)
}
