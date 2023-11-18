package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
)

func HandleResponse(w http.ResponseWriter, ctx reqcontext.RequestContext, err error, loggerMessage string, response string, httpStatus int) {
	w.WriteHeader(httpStatus)

	if err != nil {
		ctx.Logger.WithError(err).Errorf(loggerMessage)
	}

	_, err = w.Write([]byte(response))

	if err != nil {
		ctx.Logger.WithError(err).Errorf("error writing the response")
	}
}

func HandleUserValidationError(err error, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	loggerMessage := "error validating user identity"

	if err.Error() == components.BadRequestError {
		HandleResponse(w, ctx, err, loggerMessage, components.BadRequestError, http.StatusBadRequest)
	} else {
		HandleResponse(w, ctx, err, loggerMessage, components.InternalServerError, http.StatusInternalServerError)
	}
}
