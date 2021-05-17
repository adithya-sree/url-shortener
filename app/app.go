package app

import (
	"net/http"

	"github.com/adithya-sree/url-shortener/app/db"
	"github.com/adithya-sree/url-shortener/app/handler"
	"github.com/adithya-sree/url-shortener/config"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type App struct {
	log  *logrus.Logger
	mux  *chi.Mux
	conf *config.Config
}

func (app *App) Run() error {
	app.log.Infof("listening on port [%s]", app.conf.App.Port)
	return http.ListenAndServe(":"+app.conf.App.Port, app.mux)
}

func NewApp(log *logrus.Logger, conf *config.Config, redis *redis.Client) *App {
	redisWrapper := &db.RedisWrapper{
		Redis: redis,
		Log:   log,
		Conf:  conf,
	}
	handler := &handler.Handler{
		Log:   log,
		Conf:  conf,
		Redis: redisWrapper,
	}
	mux := newMux(handler)

	return &App{
		mux:  mux,
		log:  log,
		conf: conf,
	}
}

func newMux(handler *handler.Handler) *chi.Mux {
	mux := chi.NewMux()

	mux.Group(func(r chi.Router) {
		r.Get("/ecv", handler.Ecv())

		r.Group(func(r chi.Router) {
			r.Use(handler.SessionGen)
			r.Get("/{short}", handler.Redirect())

			r.Group(func(r chi.Router) {
				r.Use(handler.ValidateCreateRedirect)
				r.Post("/shortener", handler.CreateRedirect())
			})
		})
	})

	return mux
}
