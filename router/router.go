package router

import (
	"db-forum/handler"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func handleCreateOptions(ctx *fasthttp.RequestCtx) {
	options := ctx.UserValue("options").(string)
	if options == "/create" {
		handler.CreateForum(ctx)
		return
	}

	options = options[1 : len(options)-7]
	handler.CreateThread(ctx, options)
}

func CreateRouter() *fasthttprouter.Router {
	r := fasthttprouter.New()

	r.POST("/api/forum/*options", handleCreateOptions)

	r.GET("/api/forum/:slug/details", handler.GetForumInfo)
	r.GET("/api/forum/:slug/threads", handler.GetForumThreads)
	r.GET("/api/forum/:slug/users", handler.GetForumUsers)

	r.GET("/api/post/:slug/details", handler.GetPostDetails)
	r.POST("/api/post/:slug/details", handler.UpdatePost)

	r.GET("/api/service/status", handler.GetServiceStatus)
	r.POST("/api/service/clear", handler.ClearService)

	r.POST("/api/thread/:slug/create", handler.CreatePost)

	r.GET("/api/thread/:slug/details", handler.GetThreadInfo)
	r.GET("/api/thread/:slug/posts", handler.GetPost)
	r.POST("/api/thread/:slug/vote", handler.VoteForThread)

	r.GET("/api/thread/:slug", handler.GetThreadInfo)
	r.POST("/api/thread/:slug/details", handler.UpdateThread)

	r.POST("/api/user/:nickname/create", handler.CreateUser)
	r.GET("/api/user/:nickname/profile", handler.GetUser)
	r.POST("/api/user/:nickname/profile", handler.UpdateUser)

	return r
}
