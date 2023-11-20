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
	rt.router.GET("/users/", rt.wrap(rt.searchUsers))
	rt.router.PUT("/users/:username/set_username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:username/profile", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:username/stream", rt.wrap(rt.getMyStream))

	// Bans routes
	rt.router.GET("/users/:username/bans/", rt.wrap(rt.getBans))
	rt.router.PUT("/users/:username/bans/:banned_username", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/bans/:banned_username", rt.wrap(rt.unbanUser))

	// Posts routes
	rt.router.POST("/users/:username/profile/posts", rt.wrap(rt.uploadPhoto))
	rt.router.GET("/users/:username/profile/posts", rt.wrap(rt.getUserPosts))
	rt.router.GET("/users/:username/profile/:post_id", rt.wrap(rt.getPost))
	rt.router.DELETE("/users/:username/profile/:post_id", rt.wrap(rt.deletePhoto))

	return rt.router
}
