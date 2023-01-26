package controller

import (
	"forum/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	router.HandleFunc("/", h.indexPage)

	router.HandleFunc("/sign-up", h.signUp)
	router.HandleFunc("/sign-in", h.signIn)
	router.HandleFunc("/logout", h.authenticateUser(h.LogOut))

	router.HandleFunc("/create-post", h.authenticateUser(h.createPost))
	router.HandleFunc("/get-post/", h.getPost)
	router.HandleFunc("/get-posts-by-category/", h.getPostsByCategory)
	router.HandleFunc("/get-created-posts/", h.authenticateUser(h.getCreatedPost))
	router.HandleFunc("/get-liked-posts/", h.authenticateUser(h.getLikedPost))

	router.HandleFunc("/like/", h.authenticateUser(h.likePost))
	router.HandleFunc("/dislike/", h.authenticateUser(h.disLikePost))

	router.HandleFunc("/create-comment", h.authenticateUser(h.createComment))
	router.HandleFunc("/comment-like/", h.authenticateUser(h.likeComment))
	router.HandleFunc("/comment-dislike/", h.authenticateUser(h.disLikeComment))

	router.HandleFunc("/update-post", h.authenticateUser(h.updatePost))
	router.HandleFunc("/delete", h.authenticateUser(h.deletePost))

	return router
}
