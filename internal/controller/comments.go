package controller

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/service.go"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	postID, err := strconv.Atoi(r.FormValue("postid"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	author := r.FormValue("author")
	input := r.FormValue("input")

	comment := &models.Comment{
		Author: author,
		Text:   input,
		PostID: postID,
	}

	if err := h.services.CreateComment(comment); err != nil {
		if errors.Is(err, service.ErrInvalidComment) {
			h.errorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post/%d", postID), 302)
}

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	commentID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment-like/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	userRaw := r.Context().Value(ctxKeyUser)
	username := userRaw.(models.User)

	comment, err := h.services.GetCommentByID(commentID)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Comment.LikeComment(commentID, username.Username)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post/%v", comment.PostID), 302)
}

func (h *Handler) disLikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	commentID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment-dislike/"))
	if err != nil {
		h.errorPage(w, http.StatusNotFound, err.Error())
		return
	}

	userRaw := r.Context().Value(ctxKeyUser)
	username := userRaw.(models.User)

	comment, err := h.services.GetCommentByID(commentID)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Comment.DislikeComment(commentID, username.Username)
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/get-post/%v", comment.PostID), 302)
}
