package handler

import (
	"db-forum/database"
	"db-forum/models"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

func CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	if err := user.UnmarshalJSON(ctx.PostBody()); err != nil {
		WriteResponse(ctx, http.StatusBadRequest, models.Error{err.Error()})
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)
	usr, err := database.CreateUser(&user)
	if err != nil {
		if err == database.ErrDuplicate {
			WriteResponse(ctx, http.StatusConflict, usr)
			return
		}
		WriteResponse(ctx, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	WriteResponse(ctx, http.StatusCreated, (*usr)[0])
}

func GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	usr, err := database.GetUserByUsername(nickname)
	if err != nil {
		if err == database.ErrNotFound {
			WriteResponse(ctx, http.StatusNotFound, models.Error{"Can't find user\n"})
			return
		}
		log.Println(err.Error())
		WriteResponse(ctx, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	WriteResponse(ctx, http.StatusOK, usr)
}

func UpdateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	if err := user.UnmarshalJSON(ctx.PostBody()); err != nil {
		log.Println(err.Error())
		WriteResponse(ctx, http.StatusBadRequest, models.Error{err.Error()})
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)
	_, err := database.GetUserByUsername(user.Nickname)
	if err != nil {
		if err == database.ErrNotFound {
			WriteResponse(ctx, http.StatusNotFound, models.Error{"Can't find user by nickname: " + user.Nickname})
			return
		}
		WriteResponse(ctx, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	usr, err := database.UpdateUser(&user)
	if err != nil {
		if err == database.ErrDuplicate {
			WriteResponse(ctx, http.StatusConflict, models.Error{"This email is already registered by user: " + user.Nickname})
			return
		}
		if err == database.ErrNotFound {
			WriteResponse(ctx, http.StatusNotFound, models.Error{"Can't find user by nickname: " + user.Nickname})
			return
		}
		WriteResponse(ctx, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	WriteResponse(ctx, http.StatusOK, (*usr)[0])
}
