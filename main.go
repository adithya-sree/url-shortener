package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/adithya-sree/url-shortener/app"
	"github.com/adithya-sree/url-shortener/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func main() {
	log := NewLogger()
	log.Info("starting application")

	log.Info("reading configs")
	conf, err := config.ReadInConfig(config.ConfigPath)
	if err != nil {
		log.WithField(config.ErrString, err.Error()).Error("error reading config file")
		os.Exit(2)
	}

	log.Info("connecting to redis")
	redis, err := openRedis(conf, log)
	if err != nil {
		log.WithField(config.ErrString, err.Error()).Error("error connecting to redis")
		os.Exit(2)
	}

	log.Info("initializing router")
	app := app.NewApp(log, conf, redis)
	go runApp(app, log)

	log.Info("configuring safe exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("os interrupt, process exiting")
}

func runApp(app *app.App, log *logrus.Logger) {
	defer func() {
		log.Info("application no longer listening, process exiting")
		os.Exit(2)
	}()

	err := app.Run()
	if err != nil {
		log.WithField(config.ErrString, err.Error()).Error("error while listening")
	}
}

func NewLogger() *logrus.Logger {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
	}

	if strings.Contains(os.Args[0], "go-build") {
		formatter.DisableColors = false
	}

	log := logrus.New()
	log.SetFormatter(formatter)
	return log
}

func openRedis(conf *config.Config, log *logrus.Logger) (*redis.Client, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       int(conf.Redis.DB),
	}).WithTimeout(time.Duration(conf.Redis.Timeout) * time.Millisecond)
	pong, err := redis.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Infof("redis ping result [%s]", pong)
	return redis, err
}
