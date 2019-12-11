package handler

import (
	"db-forum/database"
	"db-forum/models"
	"db-forum/response"

	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/valyala/fasthttp"
	"golang.org/x/tools/container/intsets"
)

func GetForumThreads(context *fasthttp.RequestCtx) {
	slug := context.UserValue("slug").(string)
	limit := context.QueryArgs().Peek("limit")
	desc := context.QueryArgs().Peek("desc")
	since := context.QueryArgs().Peek("since")
	queryDesc, querySince := string(desc), string(since)
	var queryLimit int
	var err error
	if string(limit) == "" {
		queryLimit = intsets.MaxInt
	} else {
		queryLimit, err = strconv.Atoi(string(limit))
		if err != nil {
			response.Write(context, http.StatusBadRequest, models.Error{err.Error()})
		}
	}

	if queryDesc == "true" {
		queryDesc = "DESC"
	} else {
		queryDesc = "ASC"
	}

	_, err = database.GetForum(slug)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find forum by slug: " + slug})
			return
		}
	}

	threads, err := database.GetForumThreads(slug, querySince, queryDesc, queryLimit)
	if err != nil {
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}

	response.Write(context, http.StatusOK, (*threads))
}

func VoteForThread(context *fasthttp.RequestCtx) {
	slug := context.UserValue("slug").(string)
	body := context.PostBody()
	var voice models.Vote
	if err := voice.UnmarshalJSON(body); err != nil {
		response.Write(context, http.StatusBadRequest, models.Error{err.Error()})
		return
	}
	user, err := database.GetUserByUsername(voice.Nickname)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find user"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	voice.Nickname = user.Nickname
	var thread *models.Thread
	if govalidator.IsNumeric(slug) {
		thread, err = database.GetThread(slug, slug)
	} else {
		thread, err = database.GetThreadBySlug(slug)
	}
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find thread"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	slug = thread.Slug
	voice.ThreadId = thread.ID
	newVote, err := database.VoteThread(&voice)
	if err != nil {
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	thread.Votes = newVote
	response.Write(context, http.StatusOK, thread)
}

func UpdateThread(context *fasthttp.RequestCtx) {
	slug := context.UserValue("slug").(string)
	body := context.PostBody()
	var thread *models.Thread
	var postThread models.Thread
	if err := postThread.UnmarshalJSON(body); err != nil {
		response.Write(context, http.StatusBadRequest, models.Error{err.Error()})
		return
	}
	var err error
	if govalidator.IsNumeric(slug) {
		thread, err = database.GetThread(slug, slug)
	} else {
		thread, err = database.GetThreadBySlug(slug)
	}
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{err.Error()})
			return
		}
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	thread.Title, thread.Message = postThread.Title, postThread.Message
	resThread, err := database.UpdateThread(thread)
	if err != nil {
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusOK, resThread)
}

func GetThreadInfo(context *fasthttp.RequestCtx) {
	slug := context.UserValue("slug").(string)
	var thread *models.Thread
	var err error
	if govalidator.IsNumeric(slug) {
		thread, err = database.GetThread(slug, slug)
	} else {
		thread, err = database.GetThreadBySlug(slug)
	}

	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find forum by slug: " + slug})
			return
		}
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusOK, thread)
}

func CreateThread(context *fasthttp.RequestCtx, forumName string) {
	var thread models.Thread
	if err := thread.UnmarshalJSON(context.PostBody()); err != nil {
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	thread.Forum = forumName
	user, err := database.GetUserByUsername(thread.Author)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find user"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	thread.Author = user.Nickname
	forum, err := database.GetForum(thread.Forum)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find forum"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	thread.Forum = forum.Slug
	if thread.Slug != "" {
		existsThread, err := database.GetThreadBySlug(thread.Slug)
		if err != nil {
			if err != database.ErrNotFound {
				log.Println(err.Error())
				response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
				return
			}
		}
		if existsThread != nil {
			response.Write(context, http.StatusConflict, existsThread)
			return
		}
	}
	newThread, err := database.CreateThread(&thread)
	if err != nil {
		if err == database.ErrDuplicate {
			response.Write(context, http.StatusConflict, newThread)
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusCreated, newThread)
}
