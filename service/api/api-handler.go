package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Login routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Users routes
	rt.router.GET("/users", rt.wrap(rt.Search))

	return rt.router
}
