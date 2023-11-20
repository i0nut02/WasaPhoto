package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	var requestBody components.Photo

	err = json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		HandleResponse(w, ctx, err, "error deconding the request bidy", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	data, err := rt.db.UploadPhoto(&requestBody, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error inserting new Photo", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusCreated)
}

func (rt *_router) getUserPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	_, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retrieving username of session owner", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	usernamePosts := ps.ByName("username")

	isValid, err := rt.db.ValidUsername(usernamePosts)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to see if username in URL is valid", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	idPosts, err := rt.db.GetIdFromUsername(usernamePosts)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving URL's username id", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	isBannished, err := rt.db.IsBanished(idPosts, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to search if rewuest owner is banned", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isBannished {
		HandleResponse(w, ctx, nil, "", components.ForbiddenError, http.StatusForbidden)
		return
	}

	data, err := rt.db.GetUserPosts(idPosts, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retrieving user posts", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) getPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	_, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retrieving username of session owner", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	usernameReq := ps.ByName("username")

	idReq, err := rt.db.GetIdFromUsername(usernameReq)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving id of username URL", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	isBanned, err := rt.db.IsBanished(idReq, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking if the session owner is banned", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isBanned {
		HandleResponse(w, ctx, nil, "", components.ForbiddenError, http.StatusForbidden)
		return
	}

	post_id := ps.ByName("post_id")

	data, err := rt.db.GetPost(post_id, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to retrieve the post", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)

}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	postId := ps.ByName("post_id")

	data, err := rt.db.DeletePhoto(postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to delete the post", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusNoContent)
}
