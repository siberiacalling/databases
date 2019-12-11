package handler

import (
	"db-forum/database"
	"db-forum/models"
	"db-forum/response"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

func GetUser(context *fasthttp.RequestCtx) {
	nickname := context.UserValue("nickname").(string)
	usr, err := database.GetUserByUsername(nickname)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find user\n"})
			return
		}
		log.Println(err.Error())
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusOK, usr)
}

func CreateUser(context *fasthttp.RequestCtx) {
	var user models.User
	if err := user.UnmarshalJSON(context.PostBody()); err != nil {
		response.Write(context, http.StatusBadRequest, models.Error{err.Error()})
		return
	}
	user.Nickname = context.UserValue("nickname").(string)
	usr, err := database.CreateUser(&user)
	if err != nil {
		if err == database.ErrDuplicate {
			response.Write(context, http.StatusConflict, usr)
			return
		}
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusCreated, (*usr)[0])
}

func UpdateUser(context *fasthttp.RequestCtx) {
	var user models.User
	if err := user.UnmarshalJSON(context.PostBody()); err != nil {
		log.Println(err.Error())
		response.Write(context, http.StatusBadRequest, models.Error{err.Error()})
		return
	}
	user.Nickname = context.UserValue("nickname").(string)
	_, err := database.GetUserByUsername(user.Nickname)
	if err != nil {
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find user by nickname: " + user.Nickname})
			return
		}
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	usr, err := database.UpdateUser(&user)
	if err != nil {
		if err == database.ErrDuplicate {
			response.Write(context, http.StatusConflict, models.Error{"This email is already registered by user: " + user.Nickname})
			return
		}
		if err == database.ErrNotFound {
			response.Write(context, http.StatusNotFound, models.Error{"Can't find user by nickname: " + user.Nickname})
			return
		}
		response.Write(context, http.StatusInternalServerError, models.Error{err.Error()})
		return
	}
	response.Write(context, http.StatusOK, (*usr)[0])
}
