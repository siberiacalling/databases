package handler

import (
	"db-forum/database"
	"db-forum/models"
	"db-forum/response"
	"log"
	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"
	"golang.org/x/tools/container/intsets"
)

func CreateForum(context *fasthttp.RequestCtx) {
	body := context.PostBody()

	var forum models.Forum
	if err := forum.UnmarshalJSON(body); err != nil {
		response.Write(context, http.StatusBadRequest, models.Error{err.Error()})
		return
	}

	author, err := database.GetUserByUsername(forum.User)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find user"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	forum.User = author.Nickname
	newForum, err := database.CreateForum(&forum)
	if err != nil {
		if err == database.ErrDuplicate {
			response.Write(context, http.StatusConflict, newForum)
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, err.Error())
		return
	}
	response.Write(context, http.StatusCreated, newForum)
}

func GetForumInfo(context *fasthttp.RequestCtx) {
	id := context.UserValue("slug").(string)

	forumByID, err := database.GetForum(id)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find forumByID"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusOK, forumByID)
}

func GetForumUsers(context *fasthttp.RequestCtx) {
	id := context.UserValue("slug").(string)

	limit := string(context.QueryArgs().Peek("limit"))
	since := string(context.QueryArgs().Peek("since"))
	desc := string(context.QueryArgs().Peek("desc"))

	if limit == "" {
		limit = strconv.Itoa(intsets.MaxInt)
	}
	forumById, err := database.GetForum(id)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"can't find forumById"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	id = forumById.Slug
	users, err := database.GetForumUsers(id, limit, since, desc)
	if err != nil {
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	response.Write(context, http.StatusOK, users)
}
