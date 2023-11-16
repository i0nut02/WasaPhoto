package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user components.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil || user.Username == "" {
		w.WriteHeader(http.StatusBadRequest)

		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error parsing request body: %w", err).Error())

		_, err = w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf("error writing the response: %w", err))
		}
		return
	}
	data, err := rt.db.DoLogin(user.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error creating/searching the username: %s, details: %w", user.Username, err).Error())

		_, err = w.Write([]byte(data))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(data))

	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error writing the response, details: %w", err).Error())
	}
}
