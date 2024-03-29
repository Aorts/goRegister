package main

import (
	"goEx/config"
	register_handler "goEx/register"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err.Error())
	}
	logger := zap.Must(zap.NewProduction())

	defer logger.Sync()

	//db, err := db.InitDatabase(cfg.Database.Driver, cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Database)
	//if err != nil {
	//	panic(err.Error())
	//}
	//redisClient := redisInfra.InitRedis(cfg.Redis)
	app := fiber.New()
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.SendString("v1.1.3")
	})
	//app.Post("/api/register", register_handler.RegisterHandler(register_handler.NewRegisterFunc(db), register_handler.NewRegisterRedisFunc(redisClient)))
	//app.Post("/api/verify", register_handler.SetVerifyHandler(register_handler.NewSetVerifyFunc(redisClient), register_handler.NewDelVerifyFunc(redisClient), register_handler.NewUpdateVerifyFunc(db)))
	app.Get("/api/:cid", register_handler.GetStatusHandler(logger, register_handler.NewGetStatusFunc()))
	app.Get("/something", func(c *fiber.Ctx) error {
		return c.SendString("something")
	})

	err = app.Listen(cfg.Server.Port)
	if err != nil {
		panic(err.Error())
	}
}
