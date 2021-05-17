package handler

import (
	"net/http"

	"github.com/adithya-sree/url-shortener/app/db"
	"github.com/adithya-sree/url-shortener/config"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Log   *logrus.Logger
	Conf  *config.Config
	Redis *db.RedisWrapper
}

func (h *Handler) Ecv() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		h.Log.Info("ecv check received")
		rw.WriteHeader(http.StatusOK)
	}
}
