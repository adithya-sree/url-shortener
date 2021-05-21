package handler

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/adithya-sree/commons"
	"github.com/adithya-sree/url-shortener/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (h *Handler) ValidateCreateRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session := r.Context().Value("session").(string)
		url, err := commons.GetHeader(r, config.HeaderUrl)
		if err != nil {
			h.Log.WithField(config.Session, session).Warnf(err.Error())
			commons.RespondErrorWithSession(rw, http.StatusBadRequest, err.Error(), session)
			return
		}

		key, err := commons.GetHeader(r, config.HeaderKey)
		if err != nil {
			key = randStringBytes(h.Conf.App.ShortLength)
		}

		ctx := context.WithValue(r.Context(), "url", url)
		ctx = context.WithValue(ctx, "key", key)
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

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
