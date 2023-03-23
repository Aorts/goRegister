package main

import (
	"fmt"
	"goEx/config"
	"goEx/handler"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {

	cfg, err := initConfig()
	if err != nil {
		panic(err.Error())
	}
	db, err := initDatabase(cfg.Database.Driver, cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Database)
	if err != nil {
		panic(err.Error())
	}

	redisClient := initRedis(cfg.Redis)

	app := fiber.New()
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.SendString(cfg.Version)
	})
	app.Post("/api/register", handler.RegisterHandler(handler.NewRegisterFunc(db), handler.NewRegisterRedisFunc(redisClient)))
	app.Listen(cfg.Server.Port)
}

func initConfig() (*config.Config, error) {
	viper.SetDefault("Server.Port", ":8080")
	viper.SetDefault("Database.Version", "v0.0.1")
	viper.SetDefault("Database.Driver", "postgres")
	viper.SetDefault("Database.Username", "ts")
	viper.SetDefault("Database.Password", "ts")
	viper.SetDefault("Database.Host", "localhost:5432")
	viper.SetDefault("Database.Database", "postgres")
	viper.SetDefault("Redis.Host", "localhost:6379")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func initDatabase(driver string, username string, password string, host string, database string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v://%v:%v@%v/%v?sslmode=disable",
		driver, username, password, host, database,
	)
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Database Connected!!")
	}

	db.SetConnMaxLifetime(5 * time.Hour)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}

func initRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Host,
	})
}
