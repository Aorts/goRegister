package register_handler

import (
	"context"
	"fmt"
	"goEx/api"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type GetVerifyInput struct {
	CitizenId    string `json:"cid"`
	RegisterCode string `json:"register_code"`
}

func SetVerifyHandler(getRedisfunc SetVerifyFunc, delRedis DelVerifyFunc, updateStatus UpdateVerifyFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var veriInput GetVerifyInput

		err := c.BodyParser(&veriInput)
		if err != nil {
			return c.Status(http.StatusConflict).JSON(api.Err(409, "Conflict body"))
		}
		key := fmt.Sprintf("REGISTER:%s", veriInput.CitizenId)
		result, err := getRedisfunc(c.Context(), key)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(api.Err(404, "Invalid Citizen ID"))

		}
		if result != veriInput.RegisterCode {
			return c.Status(http.StatusNotFound).JSON(api.Err(404, "Invalid Register Code"))

		}
		err = delRedis(c.Context(), key)
		if err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(api.Err(503, "error has occurred. please contact your system administrator"))
		}
		err = updateStatus(veriInput.CitizenId)
		if err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(api.Err(503, "error has occurred. please contact your system administrator"))
		}
		return c.Status(http.StatusOK).JSON(api.Err(200, "success"))

	}
}

type SetVerifyFunc func(ctx context.Context, key string) (string, error)

type DelVerifyFunc func(ctx context.Context, key string) error

type UpdateVerifyFunc func(citizenId string) error

func NewSetVerifyFunc(redisClient *redis.Client) SetVerifyFunc {
	return func(ctx context.Context, key string) (string, error) {
		return redisClient.Get(ctx, key).Result()
	}
}

func NewDelVerifyFunc(redisClient *redis.Client) DelVerifyFunc {
	return func(ctx context.Context, key string) error {
		return redisClient.Del(ctx, key).Err()
	}
}

func NewUpdateVerifyFunc(db *sqlx.DB) UpdateVerifyFunc {
	return func(citizenId string) error {
		query := "update tbl_register set status = 'ACTIVE' , updated_date  = now()  where cid=$1"
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		res, err := tx.Exec(query, citizenId)
		if err != nil {
			return err
		}
		_ = res
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}
}
