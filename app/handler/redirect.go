package handler

import (
	"net/http"

	"github.com/adithya-sree/commons"
	"github.com/adithya-sree/url-shortener/config"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/reactivex/rxgo/v2"
	"github.com/sirupsen/logrus"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type CreatedRedirect struct {
	ResourceUrl string `json:"resourceUrl"`
}

func (h *Handler) Redirect() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(string)
		h.Log.WithField(config.Session, session).Info("redirect request")
		observable := rxgo.FromChannel(h.Redis.Get(chi.URLParam(r, "short")))
		for url := range observable.Observe() {
			if url.Error() {
				if url.E == redis.Nil {
					h.Log.WithField(config.Session, session).Error(config.UrlDoesNotExist)
					commons.RespondErrorWithSession(rw, http.StatusNotFound, config.UrlDoesNotExist, session)
					return
				}

				h.Log.WithFields(logrus.Fields{
					config.ErrString: url.E.Error(),
					config.Session:   session,
				}).Error(config.ErrorRedirecting)
				commons.RespondErrorWithSession(rw, http.StatusInternalServerError, config.ErrorRedirecting, session)
				return
			}

			redirectTo := url.V.(string)
			h.Log.WithFields(logrus.Fields{config.Url: redirectTo,
				config.Session: session}).Infof("redirecting request")
			http.Redirect(rw, r, redirectTo, http.StatusSeeOther)
			return
		}
	}
}

func (h *Handler) CreateRedirect() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(string)
		h.Log.WithField(config.Session, session).Info("create redirect request")
		short := r.Context().Value("key").(string)
		observable := rxgo.FromChannel(h.Redis.Set(short, r.Context().Value("url").(string)))
		for result := range observable.Observe() {
			if result.Error() {
				h.Log.WithFields(logrus.Fields{
					config.ErrString: result.E.Error(),
					config.Session:   session,
				}).Error(config.ErrorCreatingRedirect)
				commons.RespondErrorWithSession(rw, http.StatusInternalServerError, config.ErrorCreatingRedirect, session)
				return
			}

			h.Log.WithField(config.Session, session).Info("created new redirect")
			commons.RespondJSON(rw, http.StatusOK, &CreatedRedirect{ResourceUrl: short})
			return
		}
	}
}
