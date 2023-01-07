package controller

import (
	"errors"
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
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

		user := &models.User{
			Email:    email,
			Username: username,
			Password: password,
		}

		if err := h.services.Authorization.CreateUser(user); err != nil {
			log.Printf("Sign Up: Create User: %v", err)
			if errors.Is(err, service.ErrInvalidEmail) ||
				errors.Is(err, service.ErrInvalidUsername) ||
				errors.Is(err, service.ErrInvalidPassword) {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, RegisterError{
					ErrorMessage: "Invalid input data",
				})
				return
			}
			if errors.Is(err, service.ErrUserExist) {
				w.WriteHeader(http.StatusBadRequest)
				tmpl.Execute(w, RegisterError{
					ErrorMessage: "The username or email already exists",
				})
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		http.Redirect(w, r, "/sign-in", http.StatusFound)
	default:
		log.Println("Sign Up: Method not allowed")
		h.errorPage(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/login.html"))

	switch r.Method {
	case http.MethodGet:
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			tmpl.Execute(w, nil)
			return
		}

		if len(cookie.Value) != 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err.Error())
		}
	case http.MethodPost:
		email := r.FormValue("form-email")
		password := r.FormValue("form-password")

		token, expiresAt, err := h.services.Authorization.GenerateSessionToken(email, password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
		h.errorPage(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := h.services.DeleteSessionToken(cookie.Value); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
