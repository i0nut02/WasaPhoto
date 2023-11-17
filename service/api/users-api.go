package api

import (
	"fmt"
	"net/http"

	"encoding/json"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func HandleUserValidationError(err error, w http.ResponseWriter, ctx reqcontext.RequestContext) {
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

func (rt *_router) search(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := r.Header.Get("Authorization")
	username, err := rt.db.GetUsernameFromId(userId)

	if err != nil {
		HandleUserValidationError(err, w, ctx)
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

func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error vedifing existece of username")

		_, err = w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Errorf("error writing response")
		}
		return
	}

	if !validUsername {
		w.WriteHeader(http.StatusBadRequest)

		_, err = w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Errorf("error writing response")
		}
		return
	}

	banned, err := rt.db.IsBanned(user.Username, username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error checkin if session owner is banned")

		_, err = w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing the response")
		}
		return
	}

	if banned {
		w.WriteHeader(http.StatusForbidden)

		_, err = w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Errorf("error writing the response")
		}
		return
	}

	data, err := rt.db.TakeProfile(username, user.Username)
	fmt.Println("exit db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Errorf("error retriving profile information")

		_, err = w.Write([]byte(data))

		if err != nil {
			ctx.Logger.WithError(err).Errorf("error writing response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(data))

	if err != nil {
		ctx.Logger.WithError(err).Errorf("error writing response")
	}
	return
}
