package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
)

type User struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	Find(username string) (*User, error)
	Auth(username, password string) bool
}

type InfraServer struct {
	http.Handler
	db   UserRepository
	tmpl *template.Template
}

func NewInfraServer(db UserRepository) *InfraServer {
	srv := new(InfraServer)

	tmpl, err := template.ParseGlob("public/*.html")
	if err != nil || tmpl == nil {
		fmt.Errorf("problem opening %s %v", "public", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", srv.indexHandler)
		r.Post("/", srv.indexHandler)
		r.Get("/login", srv.loginGetHandler)
		r.Post("/login", srv.loginPostHandler)
	})
	srv.db = db
	srv.tmpl = tmpl
	srv.Handler = r

	return srv
}

func (s *InfraServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (s *InfraServer) loginGetHandler(w http.ResponseWriter, r *http.Request) {
	err := s.tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (s *InfraServer) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	switch r.Header.Get("Accept") {
	case "application/json":
		json.NewDecoder(r.Body).Decode(&user)
	default:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		user.Login = r.PostForm["username"][0]
		user.Password = r.PostForm["password"][0]
	}
	if !s.db.Auth(user.Login, user.Password) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func responder(w http.ResponseWriter, r *http.Request, tmpl string) {

}
