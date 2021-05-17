package handler

import (
	"context"
	"net/http"

	"github.com/adithya-sree/commons"
	"github.com/adithya-sree/url-shortener/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) ValidateCreateRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(string)
		key, err := commons.GetHeader(r, config.HeaderUrl)
		if err != nil {
			h.Log.WithField(config.Session, session).Warnf(err.Error())
			commons.RespondErrorWithSession(rw, http.StatusBadRequest, err.Error(), session)
			return
		}

		ctx := context.WithValue(r.Context(), "url", key)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func (h *Handler) SessionGen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, err := commons.GetHeader(r, config.HeaderSession)
		if err != nil {
			session = uuid.NewString()
		}
		ctx := context.WithValue(r.Context(), "session", session)
		h.Log.WithFields(logrus.Fields{config.Session: session, config.Url: r.URL}).Info("created session")
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
