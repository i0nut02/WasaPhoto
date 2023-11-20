package api

import (
	"net/http"

	"encoding/json"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := r.Header.Get("Authorization")
	username, err := rt.db.GetUsernameFromId(userId)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
		return
	}

	search_username := r.URL.Query().Get("search_term")

	data, err := rt.db.SearchUserBySubString(username, search_username)

	if err != nil {
		HandleResponse(w, ctx, err, "error searching maches", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		HandleResponse(w, ctx, err, "error verifing authentication", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, err, "invalid authentication", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	var user components.User

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		HandleResponse(w, ctx, err, "error decoding the body of the request", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	validNewUsername, err := user.IsValidUsername()

	if err != nil {
		HandleResponse(w, ctx, err, "error tring to check if the new username is a valid one", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !validNewUsername {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.SetUsername(id, username, user.Username)

	if err != nil {
		loggerMessage := "error setting the new username"
		var httpStatus int
		if data == components.ConflictError {
			httpStatus = http.StatusBadRequest
		} else {
			httpStatus = http.StatusInternalServerError
		}
		HandleResponse(w, ctx, err, loggerMessage, data, httpStatus)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
		return
	}

	var user components.User

	user.Username = ps.ByName("username")

	validUsername, err := rt.db.ValidUsername(user.Username)

	if err != nil {
		HandleResponse(w, ctx, err, "error vedifing existece of username", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !validUsername {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	banned, err := rt.db.IsBanished(user.Username, username)

	if err != nil {
		HandleResponse(w, ctx, err, "error checkin if session owner is banned", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if banned {
		HandleResponse(w, ctx, nil, "", components.ForbiddenError, http.StatusForbidden)
		return
	}

	data, err := rt.db.TakeProfile(username, user.Username)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving profile information", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt _router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		HandleResponse(w, ctx, err, "error validating user identity", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	toSkip, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		HandleResponse(w, ctx, err, "error converting from value into integer", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	maxQuantity, err := strconv.Atoi(r.URL.Query().Get("max_quantity"))
	if err != nil {
		HandleResponse(w, ctx, err, "error converting max_quantity value into integer", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	data, err := rt.db.GetStream(id, toSkip, maxQuantity)

	if err != nil {
		HandleResponse(w, ctx, err, "error getting the stream", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}
