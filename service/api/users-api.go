package api

import (
	"net/http"

	"encoding/json"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) search(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := r.Header.Get("Authorization")
	username, err := rt.db.GetUsernameFromId(userId)

	if err != nil {
		ctx.Logger.WithError(err).Error("error validating user identity")

		if err.Error() == components.BadRequestError {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(components.BadRequestError))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(components.InternalServerError))
		}

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	search_username := r.URL.Query().Get("search_term")

	data, err := rt.db.SearchUserBySubString(username, search_username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Errorf("error searching maches to \"%s\"", search_username)

		_, err = w.Write([]byte(data))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing the response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(data))

	if err != nil {
		ctx.Logger.WithError(err).Error("error writing the response")
	}
}

func (rt *_router) setUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	username := ps.ByName("username")

	isValid, err := rt.db.IsValid(id, username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error verifing authentication")

		_, err = w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.Error("invalid authentication")

		_, err = w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var user components.User

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error decoding the body of the request")

		_, err = w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing the response")
		}
		return
	}

	data, err := rt.db.SetUsername(id, username, user.Username)

	if err != nil {
		if data == components.ConflictError {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		ctx.Logger.WithError(err).Errorf("error setting the new username")

		_, err = w.Write([]byte(data))

		if err != nil {
			ctx.Logger.WithError(err).Errorf("error writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(data))

	if err != nil {
		ctx.Logger.WithError(err).Errorf("error writing the response")
	}
}
