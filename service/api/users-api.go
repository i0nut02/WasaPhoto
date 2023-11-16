package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := r.Header.Get("Authorization")
	username := r.Header.Get("username")

	is_valid, err := rt.db.IsValid(userId, username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error validating user identity")

		_, err = w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	if !is_valid {
		w.WriteHeader(http.StatusUnauthorized)

		_, err = w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing resonse")
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
