package controller

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/service"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type index struct {
	User     models.User
	Post     *models.Post
	Comments []*models.Comment
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/create-post.html")
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	userRaw := r.Context().Value(ctxKeyUser)
	user := userRaw.(models.User)

	switch r.Method {
	case http.MethodGet:
		post := &models.Post{}

		index := &index{
			User: user,
			Post: post,
		}

		if err = tmpl.Execute(w, index); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		r.ParseForm()
		title := r.FormValue("title")
		content := r.FormValue("content")
		about := r.FormValue("about")
		categoryString := r.Form["category"]

		post := &models.Post{
			UserID:   user.ID,
			Title:    title,
			Content:  content,
			About:    about,
			Category: categoryString,
		}

		if err = h.services.PostItem.CreatePost(post); err != nil {
			if errors.Is(err, service.ErrInvalidPost) {
				h.errorPage(w, http.StatusBadRequest, err.Error())
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) getPostsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	tmpl := template.Must(template.ParseFiles("web/template/index.html"))

	user, err := h.services.Authorization.GetSessionTokenFromRequest(r)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err.Error())
		return
	}

	category := r.URL.Query().Get("category")

	posts, err := h.services.PostItem.GetPostsByCategory(category)
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
		return
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	tmpl, err := template.ParseFiles("web/template/get-post.html")
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.Authorization.GetSessionTokenFromRequest(r)
	if err != nil {
		h.errorPage(w, http.StatusUnauthorized, err.Error())
		return
	}

	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/get-post/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	post, err := h.services.PostItem.GetPostByID(postID)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	comments, err := h.services.Comment.GetComments(postID)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
	index := &index{
		User:     user,
		Post:     &post,
		Comments: comments,
	}

	if err = tmpl.Execute(w, index); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) getCreatedPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userRaw := r.Context().Value(ctxKeyUser)
	user := userRaw.(models.User)

	posts, err := h.services.PostItem.GetCreatedPosts(user.ID)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}

	index := &Index{
		User: user,
		Post: posts,
	}

	tmpl := template.Must(template.ParseFiles("web/template/index.html"))
	if err = tmpl.Execute(w, index); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) getLikedPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	userRaw := r.Context().Value(ctxKeyUser)
	user := userRaw.(models.User)

	posts, err := h.services.PostItem.GetLikedPosts(user.Username)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	index := &Index{
		User: user,
		Post: posts,
	}

	tmpl := template.Must(template.ParseFiles("web/template/index.html"))
	err = tmpl.Execute(w, index)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/like/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	username := r.FormValue("username")

	if err = h.services.LikePost(username, id); err != nil {
		log.Println(err)
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post/%v", id), http.StatusFound)
}

func (h *Handler) disLikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/dislike/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	username := r.FormValue("username")

	if err = h.services.DisLikePost(username, id); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post/%v", id), http.StatusFound)
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/partials/layout.html", "web/template/editpost.html", "web/template/privat.navbar.html")
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	post, err := h.services.PostItem.GetPostByID(id)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	index := &index{
		Post: &post,
	}

	if r.Method == "GET" {
		tmpl.Execute(w, index)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	err = h.services.PostItem.UpdatePost(id, post.Like, post.DisLike, title, content)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	if err = h.services.DeletePost(id); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
