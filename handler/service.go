package handler

import (
	"net/http"

	"db-forum/database"

	"github.com/valyala/fasthttp"
)

func ClearService(ctx *fasthttp.RequestCtx) {
	database.ClearTable()
	WriteResponse(ctx, http.StatusOK, nil)
}

func GetServiceStatus(ctx *fasthttp.RequestCtx) {
	WriteResponse(ctx, http.StatusOK, database.GetStatus())
}
