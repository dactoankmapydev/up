package conn

/*import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisDB struct {
	Client *redis.Client
	Host   string
	Port   string
}

func (rd *RedisDB) NewRedisDB() {
	rd.Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", rd.Host, rd.Port),
	})
	_, err := rd.Client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Info(err)
	}
	logrus.Info("successful connection to redis")
}*/

