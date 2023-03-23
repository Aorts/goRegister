package handler

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type RegisterInput struct {
	CitizenId string `json:"cid"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Mobile    string `json:"mobile"`
}

func RegisterHandler(registerFunc RegisterFunc, setRegisterRedisFunc SetRegisterRedisFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regInput RegisterInput

		err := c.BodyParser(&regInput)
		if err != nil {
			return fiber.NewError(999, "cannot Parser to body")
		}
		key := fmt.Sprintf("REGISTER:%v", regInput.CitizenId)

		resgfisResult, err := registerFunc(regInput.CitizenId, regInput.Name, regInput.Birthdate, regInput.Mobile)
		if err != nil {
			if strings.Contains(err.Error(), "tbl_register_cid_key") {
				return c.JSON(fiber.Map{
					"code":    200,
					"message": "User already registerd",
				})
			} else {
				return c.JSON(fiber.Map{
					"code":    200,
					"message": "error has occurred. please contact your system administrator",
				})
			}
		}

		err = setRegisterRedisFunc(c.Context(), key, resgfisResult)
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    200,
				"message": "error has occurred. please contact your system administrator",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data": fiber.Map{
				"register_code": resgfisResult,
			},
		})
	}
}

type RegisterFunc func(citizenId string, name string, birthdate string, mobile string) (string, error)

type SetRegisterRedisFunc func(ctx context.Context, key string, value string) error

type GetRedisFunc func(ctx context.Context, key string) (string, error)

func NewRegisterFunc(db *sqlx.DB) RegisterFunc {
	return func(citizenId string, name string, birthdate string, mobile string) (string, error) {
		randNumStr := getRegisterCode()
		query := "insert into tbl_register (cid, name, birthdate, mobile, status, register_code) values ($1, $2, $3, $4, $5, $6)"

		tx, err := db.Begin()
		if err != nil {
			return "", err
		}
		res, err := tx.Exec(
			query,
			citizenId,
			name,
			birthdate,
			mobile,
			"Pending",
			randNumStr,
		)
		if err != nil {
			return "", err
		}
		err = tx.Commit()
		if err != nil {
			return "Commit Failure", err
		}
		_ = res
		return randNumStr, nil
	}
}

func getRegisterCode() string {
	// random number for register_code
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(9000) + 1000
	randNumStr := strconv.Itoa(randNum)
	return randNumStr
}

func NewRegisterRedisFunc(redisClient *redis.Client) SetRegisterRedisFunc {
	return func(ctx context.Context, key string, value string) error {
		return redisClient.Set(ctx, key, value, 0).Err()
	}
}
