package bootstrap

import (
	"go-server/mongo"

	"github.com/go-redis/redis/v8"
)

type Application struct {
	Env   *Env
	Mongo mongo.Client
	Redis *redis.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	app.Redis = NewRedisClient(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
