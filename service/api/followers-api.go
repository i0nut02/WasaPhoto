package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameId, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		if err.Error() != components.BadRequestError {
			HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		} else {
			HandleResponse(w, ctx, err, "error retriving username of auth id", components.InternalServerError, http.StatusInternalServerError)
		}
		return
	}

	usernameUrl := ps.ByName("username")

	isBanned, err := rt.db.IsBanished(usernameUrl, usernameId)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking if user is banished", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isBanned {
		HandleResponse(w, ctx, nil, "", components.ForbiddenError, http.StatusForbidden)
		return
	}

	data, err := rt.db.GetFollowers(usernameUrl, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving follower list", data, http.StatusInternalServerError)
		return
	}
	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) getFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameId, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		if err.Error() != components.BadRequestError {
			HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		} else {
			HandleResponse(w, ctx, err, "error retriving username of auth id", components.InternalServerError, http.StatusInternalServerError)
		}
		return
	}

	usernameUrl := ps.ByName("username")

	isBanned, err := rt.db.IsBanished(usernameUrl, usernameId)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking if user is banished", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isBanned {
		HandleResponse(w, ctx, nil, "", components.ForbiddenError, http.StatusForbidden)
		return
	}

	data, err := rt.db.GetFollowing(usernameUrl, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving follower list", data, http.StatusInternalServerError)
		return
	}
	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking if the session owner access it username", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	followedUsername := ps.ByName("followed_username")

	followedId, err := rt.db.GetIdFromUsername(followedUsername)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving id of the followed user", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if followedId == id {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.FollowUser(id, followedId)

	if err != nil {
		HandleResponse(w, ctx, err, "error following the followedId", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking if the session owner access it username", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	followedUsername := ps.ByName("followed_username")

	followedId, err := rt.db.GetIdFromUsername(followedUsername)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving id of the followed user", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if followedId == id {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.UnfollowUser(id, followedId)

	if err != nil {
		HandleResponse(w, ctx, err, "error following the followedId", data, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}