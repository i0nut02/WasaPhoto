package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user components.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		HandleResponse(w, ctx, err, "error parsing request body", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	isValidUsername, err := user.IsValidUsername()

	if err != nil {
		HandleResponse(w, ctx, err, "error checking if the username is valid", components.InternalServerError, http.StatusInternalServerError)
		return
	}

	if !isValidUsername {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.DoLogin(user.Username)

	if err != nil {
		HandleResponse(w, ctx, err, "error creating/searching the username", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusCreated)
}
