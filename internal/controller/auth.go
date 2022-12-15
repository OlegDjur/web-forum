package controller

import (
	"errors"
	"forum/internal/models"
	"log"
	"net/http"
	"text/template"
	"time"

	"forum/internal/service.go"
)

type RegisterError struct {
	ErrorMessage string
}

type LoginError struct {
	ErrorMessage string
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	log.Println("/register", r.Method)
	tmpl := template.Must(template.ParseFiles("web/template/registration.html"))

	switch r.Method {
	case http.MethodGet:
		if err := tmpl.Execute(w, nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		username := r.FormValue("form-username")
		email := r.FormValue("form-email")
		password := r.FormValue("form-password")

		if len(password) < 8 {
			tmpl.Execute(w, RegisterError{
				ErrorMessage: "Password must be at least 8 characters",
			})
			return
		}

		user := &models.User{
			Email:    email,
			Username: username,
			Password: password,
		}

		err := h.services.Authorization.CreateUser(user)
		if errors.Is(err, service.ErrInvalidEmail) ||
			errors.Is(err, service.ErrInvalidUsername) {
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, RegisterError{
				ErrorMessage: "The username or email already exists",
			})
			return
		}

		http.Redirect(w, r, "/sign-in", http.StatusFound)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	log.Println("/login", r.Method)
	tmpl := template.Must(template.ParseFiles("web/template/login.html"))

	if r.Method == "GET" {
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != "POST" {
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := tmpl.Execute(w, nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		email := r.FormValue("form-email")
		password := r.FormValue("form-password")

		token, expiresAt, err := h.services.Authorization.GenerateSessionToken(email, password)
		if err != nil {
			tmpl.Execute(w, LoginError{
				ErrorMessage: "Invalid email or password",
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "sessionID",
			Value:   token,
			Expires: expiresAt,
		})

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

}

func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Fatal(err)
	}
	if err := h.services.DeleteSessionToken(cookie.Value); err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
