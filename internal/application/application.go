package application

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/alextonkonogov/atonko-authorization/internal/repository"
)

var flag bool

type app struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
	repo   *repository.Repository
}

func (a app) Routes(r *httprouter.Router) {
	r.ServeFiles("/public/*filepath", http.Dir("public"))
	r.GET("/", a.authorized(a.StartPage))
	r.GET("/login", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.LoginPage(rw, "")
	})
	r.POST("/login", a.Login)
}

func (a app) authorized(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !flag {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}

		next(rw, r, ps)
	}
}

func (a app) LoginPage(rw http.ResponseWriter, message string) {
	lp := filepath.Join("public", "html", "login.html")

	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	type answer struct {
		Message string
	}
	data := answer{message}

	err = tmpl.ExecuteTemplate(rw, "login", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		a.LoginPage(rw, "Необходимо указать логин и пароль!")
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	_, err := a.repo.Login(a.ctx, a.dbpool, login, hashedPass)
	if err != nil {
		a.LoginPage(rw, "Вы ввели неверный логин или пароль!")
		return
	}

	flag = true
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) StartPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	flag = false
	motivation, err := a.repo.GetRandomMotivation(a.ctx, a.dbpool)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "motivation.html")

	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(rw, "motivation", motivation)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, dbpool, repository.NewRepository(dbpool)}
}
