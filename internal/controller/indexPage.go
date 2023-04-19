package controller

import (
	"forum/internal/models"
	"net/http"
	"text/template"
)

type Index struct {
	User models.User
	Post []models.Post
}

func (h *Handler) indexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	tmpl := template.Must(template.ParseFiles("web/template/index.html"))

	user, err := h.services.Authorization.GetSessionTokenFromRequest(r)
	// if err != nil {
	// 	h.errorPage(w, http.StatusUnauthorized, err.Error())
	// 	return
	// }

	posts, err := h.services.PostItem.GetAllPosts()
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	index := &Index{
		User: user,
		Post: posts,
	}

	if err = tmpl.Execute(w, index); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}
