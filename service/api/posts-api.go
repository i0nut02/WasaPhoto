package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) checkUserAccesAndPostAuthor(userId string, requestedUsername string, requestedPostId string) (data string, httpResponse int, err error) {
	usernameId, err := rt.db.GetUsernameFromId(userId)

	if err != nil {
		if err.Error() != components.BadRequestError {
			return components.InternalServerError, http.StatusInternalServerError, err
		} else {
			return components.BadRequestError, http.StatusBadRequest, nil
		}
	}

	_, err = rt.db.GetIdFromUsername(requestedUsername)

	if err != nil {
		if err.Error() != components.BadRequestError {
			return components.InternalServerError, http.StatusInternalServerError, err
		} else {
			return components.BadRequestError, http.StatusBadRequest, nil
		}
	}

	isBanned, err := rt.db.IsBanished(requestedUsername, usernameId)

	if err != nil {
		return components.InternalServerError, http.StatusInternalServerError, err
	}

	if isBanned {
		return components.ForbiddenError, http.StatusForbidden, nil
	}

	isAuthor, err := rt.db.ValidPostAuthor(requestedPostId, requestedUsername)

	if err != nil {
		return components.InternalServerError, http.StatusInternalServerError, err
	}

	if !isAuthor {
		return components.BadRequestError, http.StatusBadRequest, nil
	}

	return "", -1, nil
}

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

	var photo components.Photo

	err = json.NewDecoder(r.Body).Decode(&photo)

	if err != nil {
		HandleResponse(w, ctx, err, "error deconding the request bidy", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	data, err := rt.db.UploadPhoto(&photo, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error inserting new Photo", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusCreated)
}

func (rt *_router) getUserPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameId, err := rt.db.GetUsernameFromId(id)

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

	isBannished, err := rt.db.IsBanished(usernamePosts, usernameId)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to search if request owner is banned", components.InternalServerError, http.StatusInternalServerError)
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

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	data, err := rt.db.GetPost(postId, id)

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

	validAuthor, err := rt.db.ValidPostAuthor(postId, username)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking author of the post", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !validAuthor {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.DeletePhoto(postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to delete the post", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusNoContent)
}

func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	data, err := rt.db.GetLikes(postId, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving likers list", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	likerId := ps.ByName("liker_id")

	if likerId != id {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.LikePhoto(postId, likerId)

	if err != nil {
		HandleResponse(w, ctx, err, "error inserting the new like", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	likerId := ps.ByName("liker_id")

	if likerId != id {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.UnlikePhoto(postId, likerId)

	if err != nil {
		HandleResponse(w, ctx, err, "error inserting the new like", data, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	data, err := rt.db.GetComments(postId, id)

	if err != nil {
		HandleResponse(w, ctx, err, "error retrieving the comments", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	var comment components.CommentContent

	err = json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		HandleResponse(w, ctx, err, "error retriving body content", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	data, err := rt.db.CommentPhoto(postId, id, comment.Content)

	if err != nil {
		HandleResponse(w, ctx, err, "error creating comment", data, http.StatusInternalServerError)
		return
	}
	HandleResponse(w, ctx, nil, "", data, http.StatusCreated)
}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameReq := ps.ByName("username")

	postId := ps.ByName("post_id")

	response, httpResponse, err := rt.checkUserAccesAndPostAuthor(id, usernameReq, postId)

	if err != nil {
		HandleResponse(w, ctx, err, "error managing Url and session permissions", response, httpResponse)
		return
	}

	if httpResponse != -1 {
		HandleResponse(w, ctx, nil, "", response, httpResponse)
		return
	}

	isValid, err := rt.db.IsValid(id, usernameReq)

	if err != nil {
		HandleResponse(w, ctx, err, "error checking validity of username and Auth", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValid {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	commentId := ps.ByName("comment_id")

	data, err := rt.db.UncommentPhoto(commentId)

	if err != nil {
		HandleResponse(w, ctx, err, "error deleting the comment", data, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
