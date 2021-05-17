package db

import (
	"context"
	"time"

	"github.com/adithya-sree/url-shortener/config"
	"github.com/go-redis/redis/v8"
	"github.com/reactivex/rxgo/v2"
	"github.com/sirupsen/logrus"
)

type RedisWrapper struct {
	Redis *redis.Client
	Log   *logrus.Logger
	Conf  *config.Config
}

func (r *RedisWrapper) Get(key string) chan rxgo.Item {
	ch := make(chan rxgo.Item, 1)

	go func() {
		ctx, cancel := r.getContext()
		defer cancel()
		val, err := r.Redis.Get(ctx, key).Result()
		if err != nil {
			ch <- rxgo.Item{
				V: "",
				E: err,
			}

			return
		}

		ch <- rxgo.Item{
			V: val,
			E: nil,
		}
	}()

	return ch
}

func (r *RedisWrapper) Set(key, val string) chan rxgo.Item {
	ch := make(chan rxgo.Item, 1)

	go func() {
		ctx, cancel := r.getContext()
		defer cancel()
		err := r.Redis.Set(ctx, key, val, 0).Err()
		if err != nil {
			ch <- rxgo.Item{
				V: "",
				E: err,
			}

			return
		}

		ch <- rxgo.Item{
			V: "",
			E: nil,
		}

		return
	}()

	return ch
}

func (r *RedisWrapper) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(r.Conf.Redis.Timeout)*time.Microsecond)
}
