package application

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"

	"github.com/alextonkonogov/atonko-authorization/internal/repository"
)

type app struct {
	ctx    context.Context
	dbpool *pgxpool.Pool
	repo   *repository.Repository
}

func (a app) Routes(r *httprouter.Router) {
	r.ServeFiles("/public/*filepath", http.Dir("public"))
	r.GET("/", a.StartPage)
}

func (a app) StartPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
