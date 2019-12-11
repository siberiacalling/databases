package handler

import (
	"net/http"

	"db-forum/database"
	"db-forum/response"
	"github.com/valyala/fasthttp"
)

func ClearService(context *fasthttp.RequestCtx) {
	database.Clear()
	response.Write(context, http.StatusOK, nil)
}

func GetServiceStatus(context *fasthttp.RequestCtx) {
	response.Write(context, http.StatusOK, database.GetStatus())
}
