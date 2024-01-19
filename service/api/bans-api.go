package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameId, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
	}

	usernameURL := ps.ByName("username")

	if usernameId != usernameURL {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	usernameToBan := ps.ByName("banned_username")

	idToBan, err := rt.db.GetIdFromUsername(usernameToBan)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
		return
	}

	if usernameToBan == usernameId {
		HandleResponse(w, ctx, nil, "", components.BadRequestError, http.StatusBadRequest)
		return
	}

	data, err := rt.db.BanUser(id, idToBan)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to ban user", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}

func (rt _router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.Header.Get("Authorization")

	usernameId, err := rt.db.GetUsernameFromId(id)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
	}

	usernameURL := ps.ByName("username")

	if usernameId != usernameURL {
		HandleResponse(w, ctx, nil, "", components.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	usernameToBan := ps.ByName("banned_username")

	idToBan, err := rt.db.GetIdFromUsername(usernameToBan)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
		return
	}

	data, err := rt.db.UnbanUser(id, idToBan)

	if err != nil {
		HandleResponse(w, ctx, err, "error trying to ban user", data, http.StatusInternalServerError)
		return
	}

	HandleResponse(w, ctx, nil, "", data, http.StatusOK)
}
